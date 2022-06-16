package kazuha

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

func (c *char) Skill(p map[string]int) action.ActionInfo {
	hold := p["hold"]
	c.a1Ele = attributes.NoElement
	if hold == 0 {
		return c.skillPress(p)
	}
	return c.skillHold(p)
}

func (c *char) skillPress(p map[string]int) (int, int) {
	f, a := c.ActionFrames(action.ActionSkill, p)
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
	c.Core.QueueAttack(ai, combat.NewDefCircHit(1.5, false, combat.TargettableEnemy), 0, f)

	c.Core.QueueParticle("kazuha", 3, attributes.Anemo, 100)

	c.Core.Tasks.Add(c.absorbCheckA1(c.Core.F, 0, int(f/6)), "kaz-a1-absorb-check", 1)

	cd := 360
	if c.Base.Cons > 0 {
		cd = 324
	}
	if c.Base.Cons == 6 {
		c.c6Active = c.Core.F + f + 300
		c.AddWeaponInfuse(core.WeaponInfusion{
			Key:    "kazuha-c6-infusion",
			Ele:    attributes.Anemo,
			Tags:   []combat.AttackTag{combat.AttackTagNormal, combat.AttackTagExtra, combat.AttackTagPlunge},
			Expiry: c.Core.F + f + 300,
		})
	}

	c.SetCD(action.ActionSkill, cd)

	return f, a
}

func (c *char) skillHold(p map[string]int) (int, int) {
	f, a := c.ActionFrames(action.ActionSkill, p)
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

	c.Core.QueueAttack(ai, combat.NewDefCircHit(1.5, false, combat.TargettableEnemy), 0, f)

	c.Core.QueueParticle("kazuha", 4, attributes.Anemo, 100)

	c.Core.Tasks.Add(c.absorbCheckA1(c.Core.F, 0, int(f/6)), "kaz-a1-absorb-check", 1)
	cd := 540
	if c.Base.Cons > 0 {
		cd = 486
	}

	if c.Base.Cons == 6 {
		c.c6Active = c.Core.F + f + 300
		c.AddWeaponInfuse(core.WeaponInfusion{
			Key:    "kazuha-c6-infusion",
			Ele:    attributes.Anemo,
			Tags:   []combat.AttackTag{combat.AttackTagNormal, combat.AttackTagExtra, combat.AttackTagPlunge},
			Expiry: c.Core.F + f + 300,
		})
	}

	c.SetCD(action.ActionSkill, cd)
	return f, a
}
