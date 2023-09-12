package lynette

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var (
	skillPressFrames   []int
	skillHoldEndFrames []int
)

const (
	skillCD = 12 * 60

	skillPressHitmark        = 28
	skillPressAlignedHitmark = 58
	skillPressC6Start        = 17
	skillPressCDStart        = 26

	skillHoldShadowsignStart    = 17
	skillHoldShadowsignInterval = 0.1 * 60
	skillHoldEndHitmark         = 16
	skillHoldEndAlignedHitmark  = 44
	skillHoldEndC6Start         = 14 - 9 // 9f before cd start
	skillHoldEndCDStart         = 14

	particleICDKey = "lynette-particle-icd"
	particleICD    = 0.6 * 60
	particleCount  = 4

	skillTag           = "lynette-shadowsign"
	skillAlignedICDKey = "lynette-aligned-icd"
	skillAlignedICD    = 10 * 60
)

func init() {
	// Tap E
	skillPressFrames = frames.InitAbilSlice(58) // E -> Walk
	skillPressFrames[action.ActionAttack] = 43
	skillPressFrames[action.ActionSkill] = 44
	skillPressFrames[action.ActionBurst] = 45
	skillPressFrames[action.ActionDash] = 44
	skillPressFrames[action.ActionJump] = 43
	skillPressFrames[action.ActionSwap] = 42

	// Hold E
	skillHoldEndFrames = frames.InitAbilSlice(41) // Hold E -> Walk
	skillHoldEndFrames[action.ActionAttack] = 29
	skillHoldEndFrames[action.ActionSkill] = 29
	skillHoldEndFrames[action.ActionBurst] = 29
	skillHoldEndFrames[action.ActionDash] = 28
	skillHoldEndFrames[action.ActionJump] = 31
	skillHoldEndFrames[action.ActionSwap] = 27
}

func (c *char) Skill(p map[string]int) action.Info {
	hold := p["hold"]
	if hold > 0 {
		if hold > 150 {
			hold = 150
		}
		// min duration in e state: ~35f
		// max duration in e state: ~184f
		// -> offset of 34f to 1 <= hold <= 150
		return c.skillHold(p, hold+34)
	}
	return c.skillPress(p)
}

func (c *char) skillPress(p map[string]int) action.Info {
	// press attack and aligned attack
	c.QueueCharTask(func() {
		c.Core.QueueAttack(
			c.skillAI,
			combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -0.5}, 1.8, 4.5),
			0,
			0,
			c.particleCB,
			c.makeSkillHealAndDrainCB(),
		)
		c.skillAligned(skillPressAlignedHitmark - skillPressHitmark) // TODO: unsure about when the check for aligned cd happens
	}, skillPressHitmark)

	c.QueueCharTask(c.c6, skillPressC6Start)

	c.SetCDWithDelay(action.ActionSkill, skillCD, skillPressCDStart)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionSwap],
		State:           action.SkillState,
	}
}

func (c *char) skillHold(p map[string]int, duration int) action.Info {
	// shadowsign activation
	c.QueueCharTask(func() {
		c.shadowsignSrc = c.Core.F
		c.applyShadowsign(c.Core.F)
	}, skillHoldShadowsignStart)

	// shadowsign termination, hold attack and aligned attack
	c.QueueCharTask(func() {
		c.clearShadowSign()
		c.shadowsignSrc = -1 // cancel ticks because skill is over

		c.Core.QueueAttack(
			c.skillAI,
			combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -0.5}, 1.8, 5),
			0,
			0,
			c.particleCB,
			c.makeSkillHealAndDrainCB(),
		)
		c.skillAligned(skillHoldEndAlignedHitmark - skillHoldEndHitmark) // TODO: unsure about when the check for aligned cd happens
	}, duration+skillHoldEndHitmark)

	c.QueueCharTask(c.c6, duration+skillHoldEndC6Start)

	c.SetCDWithDelay(action.ActionSkill, skillCD, duration+skillHoldEndCDStart)

	return action.Info{
		Frames:          func(next action.Action) int { return duration + skillHoldEndFrames[next] },
		AnimationLength: duration + skillHoldEndFrames[action.InvalidAction],
		CanQueueAfter:   duration + skillHoldEndFrames[action.ActionSwap],
		State:           action.SkillState,
	}
}

func (c *char) applyShadowsign(src int) func() {
	return func() {
		if src != c.shadowsignSrc {
			return
		}

		c.clearShadowSign()

		// apply shadowsign to nearest enemy
		// TODO: should actually select highest score and not closest
		enemy := c.Core.Combat.ClosestEnemyWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 8), nil)
		enemy.SetTag(skillTag, 1)

		// queue up next shadowsign application
		c.QueueCharTask(c.applyShadowsign(src), skillHoldShadowsignInterval)
	}
}

// clear shadowsign from all enemies
func (c *char) clearShadowSign() {
	for _, t := range c.Core.Combat.Enemies() {
		if e, ok := t.(*enemy.Enemy); ok {
			e.SetTag(skillTag, 0)
		}
	}
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, particleICD, true)
	c.Core.QueueParticle(c.Base.Key.String(), particleCount, attributes.Anemo, c.ParticleDelay)
}

func (c *char) skillAligned(hitmark int) {
	if c.StatusIsActive(skillAlignedICDKey) {
		return
	}
	c.AddStatus(skillAlignedICDKey, skillAlignedICD, true)

	c.Core.QueueAttack(
		c.skillAlignedAI,
		combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -0.3}, 1.2, 4.5),
		hitmark,
		hitmark,
	)
}

// When the Enigma Thrust hits an opponent, it will restore Lynette's HP based on her Max HP,
// and in the 4s afterward, she will lose a certain amount of HP per second.
func (c *char) makeSkillHealAndDrainCB() combat.AttackCBFunc {
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true

		// heal
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  c.Index,
			Message: "Enigmatic Feint",
			Src:     0.25 * c.MaxHP(),
			Bonus:   c.Stat(attributes.Heal),
		})

		// drain
		// TODO: does this really stack?
		// TODO: proper frames for interval
		c.QueueCharTask(c.skillDrain(0), 1*60)
	}
}

func (c *char) skillDrain(count int) func() {
	return func() {
		count++
		if count == 4 {
			return
		}
		c.Core.Player.Drain(player.DrainInfo{
			ActorIndex: c.Index,
			Abil:       "Enigmatic Feint",
			Amount:     0.06 * c.CurrentHP(),
		})
		c.QueueCharTask(c.skillDrain(count), 1*60)
	}
}
