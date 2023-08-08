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
	skillPressFrames     []int
	skillShortHoldFrames []int
	skillHoldFrames      []int
)

const (
	// TODO: proper frames, currently using kirara
	skillCD = 12 * 60

	skillPressHitmark = 14
	skillPressCDStart = 14

	skillHoldHitmark = 2.5*60 + 36
	skillHoldCDStart = 2.5*60 + 14

	skillMaxHoldDuration = 2.5 * 60

	particleICDKey = "lynette-particle-icd"
	particleICD    = 0.6 * 60
	particleCount  = 4

	skillTag           = "lynette-shadowsign"
	skillAlignedICDKey = "lynette-aligned-icd"
	skillAlignedICD    = 10 * 60
)

func init() {
	// TODO: proper frames, currently using kirara
	// Tap E
	skillPressFrames = frames.InitAbilSlice(38) // E -> Walk
	skillPressFrames[action.ActionAttack] = 34
	skillPressFrames[action.ActionSkill] = 34
	skillPressFrames[action.ActionBurst] = 34
	skillPressFrames[action.ActionDash] = 35
	skillPressFrames[action.ActionJump] = 35
	skillPressFrames[action.ActionSwap] = 33

	// Hold E
	skillHoldFrames = frames.InitAbilSlice(skillMaxHoldDuration + 68) // Hold E -> Walk
	skillHoldFrames[action.ActionAttack] = skillMaxHoldDuration + 59
	skillHoldFrames[action.ActionSkill] = skillMaxHoldDuration + 62
	skillHoldFrames[action.ActionBurst] = skillMaxHoldDuration + 63
	skillHoldFrames[action.ActionDash] = skillMaxHoldDuration + 62
	skillHoldFrames[action.ActionJump] = skillMaxHoldDuration + 63
	skillHoldFrames[action.ActionSwap] = skillMaxHoldDuration + 63
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	hold := p["hold"]
	// TODO: remember to make this intuitive
	if hold > 0 {
		if hold > skillMaxHoldDuration {
			hold = skillMaxHoldDuration
		}
		return c.skillHold(p, hold)
	}
	return c.skillPress(p)
}

func (c *char) skillPress(p map[string]int) action.ActionInfo {
	c.skillAttack(skillPressHitmark, false)

	c.SetCDWithDelay(action.ActionSkill, skillCD, skillPressCDStart)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionSwap], // TODO: proper frames, should be earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillHold(p map[string]int, duration int) action.ActionInfo {
	// apply shadowsign on cast, then every 0.1s
	c.shadowsignSrc = c.Core.F
	c.applyShadowsign(c.Core.F)

	// queue up task to deal damage to enemy with shadowsign
	c.QueueCharTask(c.holdAttack, skillHoldHitmark-skillMaxHoldDuration+duration)

	c.SetCDWithDelay(action.ActionSkill, skillCD, skillHoldCDStart-skillMaxHoldDuration+duration)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return skillHoldFrames[next] - skillMaxHoldDuration + duration },
		AnimationLength: skillHoldFrames[action.InvalidAction] - skillMaxHoldDuration + duration,
		CanQueueAfter:   skillHoldFrames[action.ActionAttack] - skillMaxHoldDuration + duration, // TODO: proper frames, should be earliest cancel
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

		// queue up next shadowsign application in 0.1s
		c.QueueCharTask(c.applyShadowsign(src), 0.1*60)
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

func (c *char) holdAttack() {
	c.clearShadowSign()
	c.shadowsignSrc = -1 // cancel ticks because skill is over

	c.skillAttack(0, true)
}

func (c *char) skillAttack(skillHitmark int, hold bool) {
	h := 4.5
	if hold {
		h = 5
	}
	c.Core.QueueAttack(
		c.skillAI,
		combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -0.5}, 1.8, h),
		skillHitmark,
		skillHitmark,
		c.particleCB,
		c.makeSkillHealAndDrainCB(),
	)
	// TODO: check timing
	c.QueueCharTask(c.skillAligned, skillHitmark)
	// TODO: check timing
	c.QueueCharTask(c.c6, skillHitmark)
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

func (c *char) skillAligned() {
	if c.StatusIsActive(skillAlignedICDKey) {
		return
	}
	c.AddStatus(skillAlignedICDKey, skillAlignedICD, true)

	c.Core.QueueAttack(
		c.skillAlignedAI,
		// TODO: check center/offset
		combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -0.3}, 1.2, 4.5),
		// TODO: proper frames
		0.4*60,
		0.4*60,
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
