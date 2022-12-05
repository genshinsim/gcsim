package venti

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillPressFrames []int
var skillHoldFrames []int

func init() {
	// skill (press) -> x
	skillPressFrames = frames.InitAbilSlice(98)
	skillPressFrames[action.ActionAttack] = 22
	skillPressFrames[action.ActionAim] = 22   // assumed
	skillPressFrames[action.ActionSkill] = 22 // uses burst frames
	skillPressFrames[action.ActionBurst] = 22
	skillPressFrames[action.ActionDash] = 22
	skillPressFrames[action.ActionJump] = 22

	// skill (hold) -> x
	skillHoldFrames = frames.InitAbilSlice(289)
	skillHoldFrames[action.ActionHighPlunge] = 116
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Skyward Sonnett",
		AttackTag:    combat.AttackTagElementalArt,
		ICDTag:       combat.ICDTagNone,
		ICDGroup:     combat.ICDGroupDefault,
		StrikeType:   combat.StrikeTypePierce,
		Element:      attributes.Anemo,
		Durability:   50,
		Mult:         skillPress[c.TalentLvlSkill()],
		HitWeakPoint: true,
	}

	act := action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}

	cd := 360
	cdstart := 21
	hitmark := 51
	radius := 3.0
	var count float64 = 3
	if p["hold"] != 0 {
		cd = 900
		cdstart = 34
		hitmark = 74
		radius = 6
		count = 4
		ai.Mult = skillHold[c.TalentLvlSkill()]

		act = action.ActionInfo{
			Frames:          frames.NewAbilFunc(skillHoldFrames),
			AnimationLength: skillHoldFrames[action.InvalidAction],
			CanQueueAfter:   skillHoldFrames[action.ActionHighPlunge], // earliest cancel
			State:           action.SkillState,
		}
	}

	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), radius), 0, hitmark, c.c2)
	c.Core.QueueParticle("venti", count, attributes.Anemo, hitmark+c.ParticleDelay)

	c.SetCDWithDelay(action.ActionSkill, cd, cdstart)

	return act
}
