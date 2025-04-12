package mizuki

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillFrames []int

const (
	skillHitmark                  = 2
	skillActivateDMGRadius        = 5.5
	skillActivatePoise            = 30
	skillActivateDurability       = 25
	skillCdDelay                  = 23
	skillCd                       = 15 * 60
	skillCloudPoise               = 30
	skillParticleGenerations      = 4
	dreamDrifterStateKey          = "dreamdrifter-state"
	dreamDrifterBaseDuration      = 5 * 60
	dreamDrifterDurationExtension = 2.5 * 60
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

	// Activate DMG
	ai := combat.AttackInfo{
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
		c.dreamDrifterExtensionsRemaining = 2
	}

	c.Core.QueueAttack(
		ai,
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
	if c.particleGenerationsRemaining > 0 {
		c.particleGenerationsRemaining--
		c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Anemo, c.ParticleDelay)
	}
}

func (c *char) registerSkillCallbacks() {

	// Remove the dreamDrifter state when she leaves the field
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		prev := args[0].(int)

		if prev == c.Index {

		}
		c.Core.Status.Delete(dreamDrifterStateKey)
		return false
	}, "mizuki-exit")
}
