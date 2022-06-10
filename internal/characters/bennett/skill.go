package bennett

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames [][]int
var skillHoldHitmarks = [][]int{{45, 57}, {112, 121}}

const skillPressHitmark = 16

func init() {
	skillFrames = make([][]int, 5)

	// skill (press) -> x
	skillFrames[0] = frames.InitAbilSlice(42)
	skillFrames[0][action.ActionDash] = 22
	skillFrames[0][action.ActionJump] = 23
	skillFrames[0][action.ActionSwap] = 41

	// skill (hold=1) -> x
	skillFrames[1] = frames.InitAbilSlice(98)
	skillFrames[1][action.ActionBurst] = 97
	skillFrames[1][action.ActionDash] = 65
	skillFrames[1][action.ActionJump] = 66
	skillFrames[1][action.ActionSwap] = 96

	// skill (hold=1,c4) -> x
	skillFrames[2] = frames.InitAbilSlice(107)
	skillFrames[2][action.ActionDash] = 95
	skillFrames[2][action.ActionJump] = 95
	skillFrames[2][action.ActionSwap] = 106

	// skill (hold=2) -> x
	skillFrames[3] = frames.InitAbilSlice(343)
	skillFrames[3][action.ActionSkill] = 339 // uses burst frames
	skillFrames[3][action.ActionBurst] = 339
	skillFrames[3][action.ActionDash] = 231
	skillFrames[3][action.ActionJump] = 340
	skillFrames[3][action.ActionSwap] = 337

	// skill (hold=2,a4) -> x
	skillFrames[4] = frames.InitAbilSlice(175)
	skillFrames[4][action.ActionDash] = 171
	skillFrames[4][action.ActionJump] = 174
	skillFrames[4][action.ActionSwap] = 175
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	level, ok := p["hold"]
	if !ok || level < 0 || level > 2 {
		level = 0
	}

	c4Active := false
	if p["hold_c4"] == 1 && c.Base.Cons >= 4 {
		level = 1
		c4Active = true
	}

	if level != 0 {
		return c.skillHold(level, c4Active)
	}
	return c.skillPress()
}

func (c *char) skillPress() action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Passion Overload (Press)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Pyro,
		Durability: 50,
		Mult:       skill[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(ai, combat.NewDefCircHit(0.1, false, combat.TargettableEnemy), skillPressHitmark, skillPressHitmark)

	//25 % chance of 3 orbs
	var count float64 = 2
	if c.Core.Rand.Float64() < .25 {
		count++
	}
	c.Core.QueueParticle("bennett", count, attributes.Pyro, 120)

	// a4 reduce cd by 50%
	if c.StatModIsActive("bennett-field") {
		c.SetCDWithDelay(action.ActionSkill, 300/2, 14)
	} else {
		c.SetCDWithDelay(action.ActionSkill, 300, 14)
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames[0]),
		AnimationLength: skillFrames[0][action.InvalidAction],
		CanQueueAfter:   skillFrames[0][action.ActionDash], // earliest cancel
		Post:            skillFrames[0][action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillHold(level int, c4Active bool) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Passion Overload (Level %v)", level),
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Pyro,
		Durability: 25,
	}

	for i, v := range skillHold[level-1] {
		ai.Mult = v[c.TalentLvlSkill()]
		c.Core.QueueAttack(
			ai,
			combat.NewDefCircHit(0.1, false, combat.TargettableEnemy),
			skillHoldHitmarks[level-1][i],
			skillHoldHitmarks[level-1][i],
		)
	}
	if level == 2 {
		ai.Mult = explosion[c.TalentLvlSkill()]
		c.Core.QueueAttack(ai, combat.NewDefCircHit(0.1, false, combat.TargettableEnemy), 166, 166)
	}

	//user-specified c4 variant adds an additional attack that deals 135% of the second hit
	if level == 1 && c4Active {
		ai.Mult = skillHold[level-1][1][c.TalentLvlSkill()] * 1.35
		ai.Abil = "Passion Overload (C4)"
		c.Core.QueueAttack(ai, combat.NewDefCircHit(0.1, false, combat.TargettableEnemy), 94, 94)

	}

	// TODO: particle timing??
	//Bennett Hold E is guaranteed 3 orbs
	c.Core.QueueParticle("bennett", 3, attributes.Pyro, 298)

	// FIXME: do we really need to pass index here??
	applyA4 := c.StatModIsActive("bennett-field")

	// figure out which frames to return
	// 0: skill (press) -> x
	// 1: skill (hold=1) -> x
	// 2: skill (hold=1,c4) -> x
	// 3: skill (hold=2) -> x
	// 4: skill (hold=2,a4) -> x
	idx := -1
	var cd, cdDelay int
	switch level {
	case 1:
		idx = 1
		cd = 450
		cdDelay = 43
		if c4Active {
			idx = 2
		}
	case 2:
		idx = 3
		cd = 600
		cdDelay = 110
		if applyA4 {
			idx = 4
		}
	default:
		panic("bennett skill (hold) level can only be 1 or 2")
	}

	// reduce cd by 50%
	if applyA4 {
		cd /= 2
	}
	c.SetCDWithDelay(action.ActionSkill, cd, cdDelay)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames[idx]),
		AnimationLength: skillFrames[idx][action.InvalidAction],
		CanQueueAfter:   skillFrames[idx][action.ActionDash], // earliest cancel
		Post:            skillFrames[idx][action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}
