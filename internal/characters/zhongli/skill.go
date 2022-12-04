package zhongli

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillPressFrames []int
var skillHoldFrames []int

const skillPressHimark = 24
const skillHoldHitmark = 48

func init() {
	// skill (press) -> x
	skillPressFrames = frames.InitAbilSlice(38)
	skillPressFrames[action.ActionAttack] = 37
	skillPressFrames[action.ActionBurst] = 38
	skillPressFrames[action.ActionDash] = 23
	skillPressFrames[action.ActionJump] = 23
	skillPressFrames[action.ActionSwap] = 37

	// skill (hold) -> x
	skillHoldFrames = frames.InitAbilSlice(96)
	skillHoldFrames[action.ActionDash] = 55
	skillHoldFrames[action.ActionJump] = 55
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	max, ok := p["res_count"]
	if !ok {
		max = 3
	}

	h := p["hold"]
	nostele := p["hold_nostele"] > 0
	if h > 0 || nostele {
		return c.skillHold(max, !nostele)
	}
	return c.skillPress(max)
}

func (c *char) skillPress(max int) action.ActionInfo {
	c.Core.Tasks.Add(func() {
		c.newStele(1860, max)
	}, skillPressHimark)

	c.SetCDWithDelay(action.ActionSkill, 240, 22)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillHold(max int, createStele bool) action.ActionInfo {
	//hold does dmg
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Stone Stele (Hold)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 25,
		Mult:       skillHold[c.TalentLvlSkill()],
		FlatDmg:    0.019 * c.MaxHP(),
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 10), 0, skillHoldHitmark)

	//create a stele if less than zhongli's max stele count and desired by player
	if (c.steleCount < c.maxStele) && createStele {
		c.Core.Tasks.Add(func() {
			c.newStele(1860, max) //31 seconds
		}, skillHoldHitmark)
	}

	//make a shield - enemy debuff arrows appear 3-5 frames after the damage number shows up in game
	c.Core.Tasks.Add(func() {
		c.addJadeShield()
	}, skillHoldHitmark)

	c.SetCDWithDelay(action.ActionSkill, 720, 47)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillHoldFrames),
		AnimationLength: skillHoldFrames[action.InvalidAction],
		CanQueueAfter:   skillHoldFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}
