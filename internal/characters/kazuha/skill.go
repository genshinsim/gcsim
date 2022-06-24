package kazuha

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillPressFrames []int
var skillHoldFrames []int

const skillPressHitmark = 14
const skillHoldHitmark = 34

func init() {
	// skill (press) -> x
	//85 frames to float down
	skillPressFrames = frames.InitAbilSlice(85)
	//27 frames before the start of plunge animation
	skillPressFrames[action.ActionHighPlunge] = 27

	// skill (hold) -> x
	//177 frames to float down
	skillHoldFrames = frames.InitAbilSlice(177)
	//58 frames before start of plunge animation
	skillHoldFrames[action.ActionHighPlunge] = 58
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	hold := p["hold"]
	c.a1Ele = attributes.NoElement

	// why is the same code written twice..
	if hold == 0 {
		return c.skillPress(p)
	}
	return c.skillHold(p)
}

func (c *char) skillPress(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Chihayaburu",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(1.5, false, combat.TargettableEnemy), 0, skillPressHitmark)

	c.Core.QueueParticle("kazuha", 3, attributes.Anemo, 100)

	c.Core.Tasks.Add(c.absorbCheckA1(c.Core.F, 0, int(skillPressHitmark/6)), 1)

	cd := 360
	if c.Base.Cons >= 1 {
		cd = 324
	}
	if c.Base.Cons >= 6 {
		c.c6Active = c.Core.F + skillPressHitmark + 300
		c.Core.Player.AddWeaponInfuse(
			c.Index,
			"kazuha-c6-infusion",
			attributes.Anemo,
			skillPressHitmark+300,
			true,
			combat.AttackTagNormal, combat.AttackTagExtra, combat.AttackTagPlunge,
		)
	}

	c.SetCD(action.ActionSkill, cd)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionHighPlunge], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillHold(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Chihayaburu",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 50,
		Mult:       skillHold[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(ai, combat.NewDefCircHit(1.5, false, combat.TargettableEnemy), 0, skillHoldHitmark)

	c.Core.QueueParticle("kazuha", 4, attributes.Anemo, 100)

	c.Core.Tasks.Add(c.absorbCheckA1(c.Core.F, 0, int(skillHoldHitmark/6)), 1)
	cd := 540
	if c.Base.Cons >= 1 {
		cd = 486
	}
	if c.Base.Cons >= 6 {
		c.c6Active = c.Core.F + skillHoldHitmark + 300
		c.Core.Player.AddWeaponInfuse(
			c.Index,
			"kazuha-c6-infusion",
			attributes.Anemo,
			skillHoldHitmark+300,
			true,
			combat.AttackTagNormal, combat.AttackTagExtra, combat.AttackTagPlunge,
		)
	}

	c.SetCD(action.ActionSkill, cd)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillHoldFrames),
		AnimationLength: skillHoldFrames[action.InvalidAction],
		CanQueueAfter:   skillHoldFrames[action.ActionHighPlunge], // earliest cancel
		State:           action.SkillState,
	}
}
