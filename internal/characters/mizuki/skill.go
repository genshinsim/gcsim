package mizuki

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var skillFrames []int

const (
	skillHitmark                      = 2
	skillActivateDmgName              = "Aisa Utamakura Pilgrimage"
	skillActivateDmgRadius            = 5.5
	skillActivatePoise                = 30
	skillActivateDurability           = 25
	skillCdDelay                      = 23
	skillCd                           = 15 * 60
	skillParticleGenerations          = 4
	skillParticleGenerationIcd        = 0.5 * 60
	skillParticleGenerationIcdKey     = "mizuki-particle-icd"
	cloudDmgName                      = "Dreamdrifter Continuous Attack"
	cloudPoise                        = 20
	cloudDurability                   = 25
	cloudExplosionRadius              = 4
	cloudTravelTime                   = 10
	cloudFirstHit                 int = 0.45 * 60
	cloudHitInterval              int = 0.75 * 60
	dreamDrifterStateKey              = "dreamdrifter-state"
	dreamDrifterBaseDuration          = 5 * 60
	dreamDrifterSwirlBuffKey          = "mizuki-swirl-buff"
	mizukiSwapOutKey                  = "mizuki-exit"
)

func init() {
	skillFrames = frames.InitAbilSlice(46) // E -> N1
	skillFrames[action.ActionSkill] = 50   // E -> E
	skillFrames[action.ActionBurst] = 34   // E -> Q
	skillFrames[action.ActionDash] = 29    // E -> D
	skillFrames[action.ActionSwap] = 30    // E -> Swap
}

// Weaves memories of lovely dreams, entering a Dreamdrifter state where she floats above the ground, and dealing
// 1 instance of AoE Anemo DMG to nearby opponents.
//
// Dreamdrifter
//
//   - While in the Dreamdrifter state, Yumemizuki Mizuki will continuously drift forward, dealing AoE Anemo DMG to nearby
//     opponents at regular intervals.
//
//   - During this time, Yumemizuki Mizuki can control her direction of drift, and the pick-up distance of Yumemi Style
//     Special Snacks from the Elemental Burst Anraku Secret Spring Therapy will be increased.
//
//   - Increases the Swirl DMG that nearby party members deal based on Yumemizuki Mizuki's Elemental Mastery.
//
//     Dreamdrifter will end when Mizuki leaves the field or uses her Elemental Skill again.
func (c *char) Skill(p map[string]int) (action.Info, error) {

	// if used while in dreamDrifter state, cancel the state.
	if c.StatusIsActive(dreamDrifterStateKey) {
		c.cancelDreamDrifterState()
		return action.Info{
			Frames:          frames.NewAbilFunc(skillFrames),
			AnimationLength: skillFrames[action.InvalidAction],
			CanQueueAfter:   skillFrames[action.ActionSwap], // earliest cancel is swap
			State:           action.Idle,                    // is this correct?
		}, nil
	}

	// Activation DMG
	activationAttack := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       skillActivateDmgName,
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		PoiseDMG:   skillActivatePoise,
		Element:    attributes.Anemo,
		Durability: skillActivateDurability,
		Mult:       skill[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		activationAttack,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.Player(),
			nil,
			skillActivateDmgRadius,
		),
		0,
		skillHitmark,
		c.particleCB,
	)

	c.particleGenerationsRemaining = skillParticleGenerations

	if c.Base.Ascension >= 1 {
		c.dreamDrifterExtensionsRemaining = dreamDrifterExtensions
	}

	c.applyDreamDrifterEffect()

	c.SetCDWithDelay(action.ActionSkill, skillCd, skillCdDelay)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap], // earliest cancel is swap
		State:           action.SkillState,
	}, nil
}

func (c *char) applyDreamDrifterEffect() {
	c.AddStatus(dreamDrifterStateKey, dreamDrifterBaseDuration, true)

	c.startCloudAttacks()

	if c.Base.Cons >= 1 {
		c.applyC1Effect()
	}
}

func (c *char) skillInit() {
	for _, char := range c.Core.Player.Chars() {
		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase(dreamDrifterSwirlBuffKey, -1),
			Amount: func(ai combat.AttackInfo) (float64, bool) {
				if !c.StatusIsActive(dreamDrifterStateKey) {
					return 0, false
				}
				// These flags imply AOE Swirl, in which case this Swirl DMG bonus does not apply because
				// it was calculated in a prior call of this callback. In these cases the other reaction bonuses
				// apply instead (e.g. Melt DMG Bonus, Aggravate DMG Bonus, etc.)
				if ai.Amped || ai.Catalyzed {
					return 0, false
				}
				switch ai.AttackTag {
				case attacks.AttackTagSwirlCryo:
				case attacks.AttackTagSwirlElectro:
				case attacks.AttackTagSwirlHydro:
				case attacks.AttackTagSwirlPyro:
				default:
					return 0, false
				}

				return swirlDMG[c.TalentLvlSkill()] * c.Stat(attributes.EM), false
			},
		})
	}

	// Remove the dreamDrifter state when she leaves the field
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		prev := args[0].(int)

		if prev == c.Index && c.StatusIsActive(dreamDrifterStateKey) {
			c.cancelDreamDrifterState()
		}

		return false
	}, mizukiSwapOutKey)
}

func (c *char) startCloudAttacks() {

	// clouds DMG snapshots on activation
	cloudAttack := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       cloudDmgName,
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupMizukiSkill,
		StrikeType: attacks.StrikeTypeDefault,
		PoiseDMG:   cloudPoise,
		Element:    attributes.Anemo,
		Durability: cloudDurability,
		Mult:       cloudDMG[c.TalentLvlSkill()],
	}

	snap := c.Snapshot(&cloudAttack)

	var cloudFunc func()
	cloudFunc = func() {
		if !c.StatusIsActive(dreamDrifterStateKey) {
			return
		}
		c.Core.QueueAttackWithSnap(cloudAttack, snap, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, cloudExplosionRadius), cloudTravelTime, c.particleCB)
		c.QueueCharTask(cloudFunc, cloudHitInterval)
	}

	// First cloud is launched at approximately 0.45s after skill activation.
	c.QueueCharTask(cloudFunc, cloudFirstHit)
}

// Generates up to 4 particles on each E DMG either on activation or cloud.
func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}

	if c.StatusIsActive(skillParticleGenerationIcdKey) {
		return
	}

	if c.particleGenerationsRemaining > 0 {
		c.AddStatus(skillParticleGenerationIcdKey, skillParticleGenerationIcd, false)
		c.particleGenerationsRemaining--
		c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Anemo, c.ParticleDelay)
	}
}

func (c *char) cancelDreamDrifterState() {
	c.DeleteStatus(dreamDrifterStateKey)

	c.Core.Log.NewEvent("DreamDrifter effect cancelled", glog.LogCharacterEvent, c.Index)
}
