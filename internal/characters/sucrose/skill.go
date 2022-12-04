package sucrose

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames []int

func init() {
	skillFrames = frames.InitAbilSlice(57)
	skillFrames[action.ActionCharge] = 56
	skillFrames[action.ActionSkill] = 56
	skillFrames[action.ActionDash] = 11
	skillFrames[action.ActionJump] = 11
	skillFrames[action.ActionSwap] = 56
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Astable Anemohypostasis Creation-6308",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}

	done := false
	cb := func(_ combat.AttackCB) {
		if done {
			return
		}
		done = true
		c.a4()
		c.Core.QueueParticle("sucrose", 4, attributes.Anemo, c.ParticleDelay)
	}

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 5}, 6), 0, 42, cb)

	//reduce charge by 1
	c.SetCDWithDelay(action.ActionSkill, 900, 9)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}
