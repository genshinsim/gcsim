package venti

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
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
		AttackTag:    attacks.AttackTagElementalArt,
		ICDTag:       attacks.ICDTagNone,
		ICDGroup:     combat.ICDGroupDefault,
		StrikeType:   attacks.StrikeTypePierce,
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
	trg := c.Core.Combat.PrimaryTarget()
	var count float64 = 3
	if p["hold"] != 0 {
		cd = 900
		cdstart = 34
		hitmark = 74
		radius = 6
		trg = c.Core.Combat.Player()
		count = 4
		ai.Mult = skillHold[c.TalentLvlSkill()]

		act = action.ActionInfo{
			Frames:          frames.NewAbilFunc(skillHoldFrames),
			AnimationLength: skillHoldFrames[action.InvalidAction],
			CanQueueAfter:   skillHoldFrames[action.ActionHighPlunge], // earliest cancel
			State:           action.SkillState,
		}
	}

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(trg, nil, radius), 0, hitmark, c.c2, c.makeParticleCB(count))

	c.SetCDWithDelay(action.ActionSkill, cd, cdstart)

	return act
}

func (c *char) makeParticleCB(count float64) combat.AttackCBFunc {
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true
		c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Anemo, c.ParticleDelay)
	}
}
