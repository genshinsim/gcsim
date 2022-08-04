package heizou

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var skillEndFrames []int

func init() {
	skillEndFrames = frames.InitAbilSlice(10)
	skillEndFrames[action.ActionAttack] = 19
	skillEndFrames[action.ActionBurst] = 19
	skillEndFrames[action.ActionCharge] = 19

}

const skillHitmark = 20
const skillCDStart = 18

// if you hold while at 4 stacks it takes 17 extra seconds to release
const holdAtFullStacksPenalty = 17

func (c *char) skillHoldDuration(stacks int) int {
	//animation duration only
	//diff is the number of stacks we must charge up to reach the desired state
	diff := stacks - c.decStack
	if diff < 0 {
		diff = 0
	}
	if diff > 4 {
		diff = 4
	}
	//it's .75s per stack
	return 45 * diff
}

func (c *char) addDecStack() {
	if c.decStack < 4 {
		c.decStack++
		c.Core.Log.NewEvent("declension stack gained", glog.LogCharacterEvent, c.Index).
			Write("stacks", c.decStack)
	}
}

func (c *char) skillRelease(p map[string]int, delay int) action.ActionInfo {

	c.QueueCharTask(func() {
		hitDelay := skillHitmark - skillCDStart
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Heartstopper Strike",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Anemo,
			Durability: 50,
			Mult:       skill[c.TalentLvlAttack()] + float64(c.decStack)*decBonus[c.TalentLvlAttack()],
		}
		AoE := 0.3
		if c.decStack == 4 {
			ai.Mult += convicBonus[c.TalentLvlAttack()]
			AoE = 1
			c.Core.QueueParticle("heizou", 4, attributes.Cryo, 15+c.Core.Flags.ParticleDelay)
		}

		skillCB := func(a combat.AttackCB) {
			c.decStack = 0
			c.a4()
		}

		c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), AoE, false, combat.TargettableEnemy), hitDelay, hitDelay, skillCB)
		c.SetCD(action.ActionSkill, 10*60)

		count := 2.0
		switch c.decStack {
		case 2, 3:
			if c.Core.Rand.Float64() < .5 {
				count++
			}
		case 4:
			count++
		}
		c.Core.QueueParticle("heizou", count, attributes.Anemo, hitDelay+100)

	}, skillCDStart+delay)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return delay + skillEndFrames[next] + skillHitmark },
		AnimationLength: delay + skillEndFrames[action.InvalidAction] + skillHitmark,
		CanQueueAfter:   delay + skillEndFrames[action.ActionDash] + skillHitmark, // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillHold(p map[string]int) action.ActionInfo {
	if c.decStack == 4 {
		return c.skillRelease(p, holdAtFullStacksPenalty)
	} else {
		for i := c.decStack + 1; i <= 4; i++ {
			c.QueueCharTask(c.addDecStack, c.skillHoldDuration(i))
		}
		return c.skillRelease(p, c.skillHoldDuration(4))
	}
}

func (c *char) skillPress(p map[string]int) action.ActionInfo {
	return c.skillRelease(p, 0)
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	if p["hold"] != 0 {
		return c.skillHold(p)
	}
	return c.skillPress(p)
}
