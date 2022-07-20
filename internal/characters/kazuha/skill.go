package kazuha

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillPressFrames []int
var skillHoldFrames []int

const (
	skillPressHitmark = 10
	skillPressCDStart = 8
	skillHoldHitmark  = 33
	skillHoldCDStart  = 31
)

func init() {
	//TODO: glide cancel
	// skill (press) -> x
	//85 frames to float down
	skillPressFrames = frames.InitAbilSlice(77) //averaged all abils
	//27 frames before the start of plunge animation
	skillPressFrames[action.ActionHighPlunge] = 24

	// skill (hold) -> x
	//177 frames to float down
	skillHoldFrames = frames.InitAbilSlice(175) //averaged all abils
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
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 1.5, false, combat.TargettableEnemy), 0, skillPressHitmark)

	c.Core.QueueParticle("kazuha", 3, attributes.Anemo, skillPressHitmark+c.Core.Flags.ParticleDelay)

	c.Core.Tasks.Add(c.absorbCheckA1(c.Core.F, 0, int(skillPressHitmark/6)), 1)

	cd := 360
	if c.Base.Cons >= 1 {
		cd = 324
	}
	if c.Base.Cons >= 6 {
		c.c6()
	}

	c.SetCDWithDelay(action.ActionSkill, cd, skillPressCDStart)

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

	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 1.5, false, combat.TargettableEnemy), 0, skillHoldHitmark)

	c.Core.QueueParticle("kazuha", 4, attributes.Anemo, skillHoldHitmark+c.Core.Flags.ParticleDelay)

	c.Core.Tasks.Add(c.absorbCheckA1(c.Core.F, 0, int(skillHoldHitmark/6)), 1)
	cd := 540
	if c.Base.Cons >= 1 {
		cd = 486
	}
	if c.Base.Cons >= 6 {
		c.c6()
	}

	c.SetCDWithDelay(action.ActionSkill, cd, skillHoldCDStart)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillHoldFrames),
		AnimationLength: skillHoldFrames[action.InvalidAction],
		CanQueueAfter:   skillHoldFrames[action.ActionHighPlunge], // earliest cancel
		State:           action.SkillState,
	}
}
