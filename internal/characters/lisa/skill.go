package lisa

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillHitmarks = []int{22, 117}

//p = 0 for no hold, p = 1 for hold
func (c *char) Skill(p map[string]int) action.ActionInfo {
	hold := p["hold"]
	if hold == 1 {
		return c.skillHold(p)
	}
	return c.skillPress(p)
}

//TODO: how long do stacks last?
func (c *char) skillPress(p map[string]int) (int, int) {
	f, a := c.ActionFrames(action.ActionSkill, p)
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Violet Arc",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagLisaElectro,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skillPress[c.TalentLvlSkill()],
	}

	done := false
	cb := func(a combat.AttackCB) {
		if done {
			return
		}
		count := a.Target.GetTag(conductiveTag)
		if count < 3 {
			a.Target.SetTag(conductiveTag, count+1)
		}
		done = true
	}

	c.Core.Combat.QueueAttack(ai, combat.NewDefSingleTarget(1, combat.TargettableEnemy), 0, skillHitmarks[0], cb)

	if c.Core.Rand.Float64() < 0.5 {
		c.QueueParticle("Lisa", 1, attributes.Electro, f+100)
	}

	c.SetCDWithDelay(action.ActionSkill, 60, 17)
	return f, a
}

func (c *char) skillHold(p map[string]int) (int, int) {
	f, a := c.ActionFrames(action.ActionSkill, p)
	//no multiplier as that's target dependent
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Violet Arc (Hold)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 50,
	}

	//c2 add defense? no interruptions either way
	if c.Base.Cons >= 2 {
		//increase def for the duration of this abil in however many frames
		val := make([]float64, attributes.EndStatType)
		val[attributes.DEFP] = 0.25
		c.AddStatMod("lisa-c2",

			c.Core.F+126, attributes.NoStat, func() ([]float64, bool) { return val, true })

	}

	count := 0
	var c1cb func(a combat.AttackCB)
	if c.Base.Cons > 0 {
		c1cb = func(a combat.AttackCB) {
			if count == 5 {
				return
			}
			count++
			c.AddEnergy("lisa-c1", 2)
		}
	}

	//[8:31 PM] ArchedNosi | Lisa Unleashed: yeah 4-5 50/50 with Hold
	//[9:13 PM] ArchedNosi | Lisa Unleashed: @gimmeabreak actually wait, xd i noticed i misread my sheet, Lisa Hold E always gens 5 orbs
	c.Core.Combat.QueueAttack(ai, combat.NewDefCircHit(3, false, combat.TargettableEnemy), 0, skillHitmarks[1], c1cb)

	// count := 4
	// if c.Core.Rand.Float64() < 0.5 {
	// 	count = 5
	// }
	c.QueueParticle("Lisa", 5, attributes.Electro, f+100)

	// c.CD[def.SkillCD] = c.Core.F + 960 //16seconds, starts after 114 frames
	c.SetCDWithDelay(action.ActionSkill, 960, 114)
	return f, a
}
