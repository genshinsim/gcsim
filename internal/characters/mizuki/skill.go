package mizuki

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var skillFrames []int

const (
	skillHitmark                      = 2
	skillActivateDMGRadius            = 5.5
	skillActivatePoise                = 30
	skillActivateDurability           = 25
	skillCdDelay                      = 23
	skillCd                           = 15 * 60
	skillParticleGenerations          = 4
	skillParticleGenerationIcd        = 0.5 * 60
	skillParticleGenerationIcdKey     = "mizuki-particle-icd"
	cloudPoise                        = 20
	cloudDurability                   = 25
	cloudExplosionRadius              = 4
	cloudTravelTime                   = 10
	cloudFirstHit                 int = 0.25 * 60
	cloudHitInterval              int = 0.75 * 60
	dreamDrifterStateKey              = "dreamdrifter-state"
	dreamDrifterBaseDuration          = 5 * 60
)

func init() {
	skillFrames = frames.InitAbilSlice(46) // E -> N1
	skillFrames[action.ActionSkill] = 50   // E -> E
	skillFrames[action.ActionBurst] = 34   // E -> Q
	skillFrames[action.ActionDash] = 29    // E -> D
	skillFrames[action.ActionSwap] = 30    // E -> Swap
}

// Mizuki skill. Things that happen when you tap E:
// - DMG enemies with anemo
// - Mizuki enters DreamDrifter state (lasts 5s):
//   - Team gains Swirl DMG Bonus based on Mizuki Em
//   - Mizuki starts floating fowards
//   - Projectiles attack nearby enemies every 0.75s
//   - Mizuki cannot do anything appart:
//   - Tap E again (Cancels state)
//   - Burst (does not cancel state, skill still attacks while in animation)
//   - Dash, she dashes while in state. Only usefull for dodging, affects nothing else.
//   - Swap (Cancels state)
func (c *char) Skill(p map[string]int) (action.Info, error) {

	// Activation DMG
	activationAttack := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Aisa Utamakura Pilgrimage",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		PoiseDMG:   skillActivatePoise,
		Element:    attributes.Anemo,
		Durability: skillActivateDurability,
		Mult:       skill[c.TalentLvlSkill()],
	}

	c.particleGenerationsRemaining = skillParticleGenerations

	if c.Base.Ascension >= 1 {
		c.dreamDrifterExtensionsRemaining = dreamDrifterExtensions
	}

	c.Core.QueueAttack(
		activationAttack,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.Player(),
			nil,
			skillActivateDMGRadius,
		),
		0,
		skillHitmark,
		c.particleCB,
	)

	c.SetCDWithDelay(action.ActionSkill, skillCd, skillCdDelay)

	c.Core.Status.Add(dreamDrifterStateKey, dreamDrifterBaseDuration)

	// clouds DMG snapshots on activation
	cloudAttack := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Dreamdrifter Continuous Attack",
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

	c.Core.QueueAttackWithSnap(cloudAttack, snap, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, cloudExplosionRadius), cloudFirstHit+cloudTravelTime)

	var hitFunc func()
	hitFunc = func() {
		if !c.StatusIsActive(dreamDrifterStateKey) {
			return
		}
		c.Core.QueueAttackWithSnap(cloudAttack, snap, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, cloudExplosionRadius), cloudTravelTime)
		c.QueueCharTask(hitFunc, cloudHitInterval)
	}
	c.QueueCharTask(hitFunc, cloudHitInterval)

	if c.Base.Cons >= 1 {
		var c1Func func()
		c1Func = func() {
			if !c.StatusIsActive(dreamDrifterStateKey) {
				return
			}
			for _, target := range c.Core.Combat.Enemies() {
				if e, ok := target.(*enemy.Enemy); ok {
					e.AddStatus(c1Key, c1Duration, false)
				}
			}
			c.QueueCharTask(hitFunc, c1Interval)
		}
		c.QueueCharTask(c1Func, c1Interval)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap], // earliest cancel is swap
		State:           action.SkillState,
	}, nil
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

func (c *char) registerSkillCallbacks() {

	// Remove the dreamDrifter state when she leaves the field
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		prev := args[0].(int)

		if prev == c.Index {
			c.Core.Status.Delete(dreamDrifterStateKey)
		}

		return false
	}, "mizuki-exit")
}
