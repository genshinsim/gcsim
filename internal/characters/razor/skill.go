package razor

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var skillPressFrames []int
var skillHoldFrames []int

const (
	skillPressHitmark     = 74
	skillHoldHitmark      = 92
	skillSigilDurationKey = "razor-sigil-duration"
)

func init() {
	skillPressFrames = frames.InitAbilSlice(74)
	skillHoldFrames = frames.InitAbilSlice(92)
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	if p["hold"] > 0 {
		return c.SkillHold()
	}
	return c.SkillPress()
}

func (c *char) SkillPress() action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Claw and Thunder (Press)",
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		Element:            attributes.Electro,
		Durability:         25,
		Mult:               skillPress[c.TalentLvlSkill()],
		HitlagHaltFrames:   0.1 * 60,
		HitlagFactor:       0.03,
		CanBeDefenseHalted: true,
	}

	var c4cb func(a combat.AttackCB)
	if c.Base.Cons >= 4 {
		c4cb = c.c4cb
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), 2, false, combat.TargettableEnemy),
		skillPressHitmark,
		skillPressHitmark,
		c4cb,
	)

	c.Core.Tasks.Add(c.addSigil, skillPressHitmark)

	cd := 6 * 0.82 * 60 // A1: Decreases Claw and Thunder's CD by 18%.
	c.SetCDWithDelay(action.ActionSkill, int(cd), skillPressHitmark)

	if !c.StatusIsActive(burstBuffKey) {
		//TODO: this delay used to be 80?
		c.Core.QueueParticle("razor", 3, attributes.Electro, skillPressHitmark+c.Core.Flags.ParticleDelay)
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressHitmark,
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
		combat.NewCircleHit(c.Core.Combat.Player(), 5, false, combat.TargettableEnemy),
		skillHoldHitmark,
		skillHoldHitmark,
	)

	c.Core.Tasks.Add(c.clearSigil, skillHoldHitmark)

	cd := 10 * 0.82 * 60 // A1: Decreases Claw and Thunder's CD by 18%.
	c.SetCDWithDelay(action.ActionSkill, int(cd), skillHoldHitmark)

	if !c.StatusIsActive(burstBuffKey) {
		//TODO: this delay used to be 80?
		c.Core.QueueParticle("razor", 4, attributes.Electro, skillHoldHitmark+c.Core.Flags.ParticleDelay)
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillHoldFrames),
		AnimationLength: skillHoldFrames[action.InvalidAction],
		CanQueueAfter:   skillHoldHitmark,
		State:           action.SkillState,
	}
}

func (c *char) addSigil() {
	if !c.StatusIsActive(skillSigilDurationKey) {
		c.sigils = 0
	}

	if c.sigils < 3 {
		c.sigils++
	}
	c.AddStatus(skillSigilDurationKey, 1080, true) //18 seconds
}

func (c *char) clearSigil() {
	if !c.StatusIsActive(skillSigilDurationKey) {
		c.sigils = 0
		return
	}

	if c.sigils > 0 {
		c.AddEnergy("razor", float64(c.sigils)*5)
		c.sigils = 0
		c.DeleteStatus(skillSigilDurationKey)
	}
}

func (c *char) energySigil() {
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("er-sigil", -1),
		AffectedStat: attributes.ER,
		Amount: func() ([]float64, bool) {
			if c.StatusIsActive(skillSigilDurationKey) {
				c.skillSigilBonus[attributes.ER] = float64(c.sigils) * 0.2
				return c.skillSigilBonus, true
			}
			return nil, false
		},
	})
}
