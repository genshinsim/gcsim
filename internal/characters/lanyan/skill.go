package lanyan

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

const (
	leapBackStatus = "leap-back"

	leapBackDuration   = 66
	hitGenerateShield  = 11
	missGenerateShield = 15
	detectHitmark      = 7
)

var (
	skillMissFrames []int
	skillHitFrames  []int
)

func init() {
	skillMissFrames = frames.InitAbilSlice(49) // E -> Swap
	skillMissFrames[action.ActionAttack] = 45
	skillMissFrames[action.ActionBurst] = 46
	skillMissFrames[action.ActionDash] = 47
	skillMissFrames[action.ActionJump] = 47

	skillHitFrames = frames.InitAbilSlice(80) // E -> Swap
	skillHitFrames[action.ActionAttack] = 33
	skillHitFrames[action.ActionSkill] = 33
	skillHitFrames[action.ActionBurst] = 79
	skillHitFrames[action.ActionDash] = 77
	skillHitFrames[action.ActionJump] = 79
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(leapBackStatus) {
		return c.reathermoonRings(), nil
	}

	hold, ok := p["hold"]
	switch {
	case !ok:
	case hold < 0:
		hold = 0
	case hold > 610:
		hold = 610
	}

	c.particleGenerated = false
	c.absorbedElement = attributes.Anemo
	c.DeleteStatus(leapBackStatus)

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Swallow-Wisp Pinion Dance: Detect",
		AttackTag:  attacks.AttackTagNone,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Physical,
	}
	ap := combat.NewBoxHitOnTarget(c.Core.Combat.Player(), nil, 0.1, 5) // TODO: approximate hitbox

	c.Core.QueueAttack(ai, ap, detectHitmark+hold, detectHitmark+hold, c.leapBack())

	c.Core.Tasks.Add(func() {
		if !c.hasShield() {
			c.genShield(attributes.Anemo)
		}
	}, detectHitmark+hold+missGenerateShield)

	c.SetCDWithDelay(action.ActionSkill, 16*60, 4+hold)

	return action.Info{
		Frames: func(next action.Action) int {
			skillFrames := c.getCurrentSkillFrames()
			return skillFrames[next] + hold
		},
		AnimationLength: skillHitFrames[action.InvalidAction] + hold,
		CanQueueAfter:   detectHitmark + hold, // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) getCurrentSkillFrames() []int {
	if c.StatusIsActive(leapBackStatus) {
		return skillHitFrames
	}
	return skillMissFrames
}

func (c *char) leapBack() func(combat.AttackCB) {
	done := false
	return func(a combat.AttackCB) {
		if done {
			return
		}
		e, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		done = true
		c.AddStatus(leapBackStatus, leapBackDuration, true)

		// generate the shield with element absorb
		c.absorbedElement = c.absorbA1(e)
		c.Core.Player.Tasks.Add(func() { c.genShield(c.absorbedElement) }, hitGenerateShield)
	}
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.particleGenerated {
		return
	}
	c.particleGenerated = true
	c.Core.QueueParticle(c.Base.Key.String(), 3, attributes.Anemo, c.ParticleDelay)
}
