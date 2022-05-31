package razor

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillPressFrames, skillHoldFrames []int

const (
	skillPressHitmark = 74
	skillHoldHitmark  = 92
)

func (c *char) skillPressFrameFunc(next action.Action) int {
	return skillPressFrames[next]
}

func (c *char) skillHoldFrameFunc(next action.Action) int {
	return skillHoldFrames[next]
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	if p["hold"] > 0 {
		return c.SkillHold()
	}
	return c.SkillPress()
}

func (c *char) SkillPress() action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Claw and Thunder (Press)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skillPress[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewDefCircHit(2, false, combat.TargettableEnemy),
		skillPressHitmark,
		skillPressHitmark,
		c.c4cb,
	)

	c.AddSigil()

	cd := 6 * 0.82 * 60 // A1: Decreases Claw and Thunder's CD by 18%.
	c.SetCD(action.ActionSkill, int(cd))

	if c.Core.Status.Duration("razorburst") == 0 {
		c.Core.QueueParticle("razor", 3, attributes.Electro, 80)
	}

	return action.ActionInfo{
		Frames:          c.skillPressFrameFunc,
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressHitmark,
		Post:            skillPressHitmark,
		State:           action.SkillState,
	}
}

func (c *char) SkillHold() action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Claw and Thunder (Hold)",
		AttackTag:  combat.AttackTagElementalArtHold,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skillHold[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewDefCircHit(5, false, combat.TargettableEnemy),
		skillHoldHitmark,
		skillHoldHitmark,
	)

	c.ClearSigil()

	cd := 10 * 0.82 * 60 // A1: Decreases Claw and Thunder's CD by 18%.
	c.SetCD(action.ActionSkill, int(cd))

	if c.Core.Status.Duration("razorburst") == 0 {
		c.Core.QueueParticle("razor", 4, attributes.Electro, 80)
	}

	return action.ActionInfo{
		Frames:          c.skillHoldFrameFunc,
		AnimationLength: skillHoldFrames[action.InvalidAction],
		CanQueueAfter:   skillHoldHitmark,
		Post:            skillHoldHitmark,
		State:           action.SkillState,
	}
}

func (c *char) AddSigil() {
	if c.Core.F > c.sigilsDuration {
		c.sigils = 0
	}

	if c.sigils < 3 {
		c.sigils++
		c.sigilsDuration = c.Core.F + 18*60
	}
}

func (c *char) ClearSigil() {
	if c.Core.F > c.sigilsDuration {
		c.sigils = 0
	}

	if c.sigils > 0 {
		c.AddEnergy("razor", float64(c.sigils)*5)
		c.sigils = 0
		c.sigilsDuration = 0
	}
}

func (c *char) EnergySigil() {
	val := make([]float64, attributes.EndStatType)
	c.AddStatMod("er-sigil", -1, attributes.ER, func() ([]float64, bool) {
		if c.Core.F > c.sigilsDuration {
			return nil, false
		}

		val[attributes.ER] = float64(c.sigils) * 0.2
		return val, true
	})
}
