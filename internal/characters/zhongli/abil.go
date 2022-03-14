package zhongli

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) Attack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAttack, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  coretype.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
		FlatDmg:    0.0139 * c.HPMax,
	}

	for i := 0; i < hits[c.NormalCounter]; i++ {
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, coretype.TargettableEnemy), f-i, f-i)
	}

	c.AdvanceNormalIndex()
	return f, a
}

func (c *char) ChargeAttack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionCharge, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge",
		AttackTag:  coretype.AttackTagExtra,
		ICDTag:     core.ICDTagExtraAttack,
		ICDGroup:   core.ICDGroupPole,
		Element:    core.Physical,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
		FlatDmg:    0.0139 * c.HPMax,
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, coretype.TargettableEnemy), f-1, f-1)

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {

	cd := 240
	f, a := c.ActionFrames(core.ActionSkill, p)

	max, ok := p["res_count"]
	if !ok {
		max = 3
	}

	//press does no dmg
	if p["hold"] == 1 {
		c.skillHold(f, max, true)
		cd = 720
	} else if p["hold_nostele"] == 1 {
		c.skillHold(f, max, false)
		cd = 720
	} else {
		c.skillPress(f, max)
	}

	c.SetCD(core.ActionSkill, cd)
	//no geo drain
	return f, a
}

func (c *char) skillPress(f, max int) {
	c.AddTask(func() {
		c.newStele(1860, max)
	}, "zhongli-create-stele", f)
}

func (c *char) skillHold(f, max int, createStele bool) {
	//hold does dmg
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Stone Stele (Hold)",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagElementalArt,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Geo,
		Durability: 25,
		Mult:       skillHold[c.TalentLvlSkill()],
		FlatDmg:    0.019 * c.HPMax,
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, coretype.TargettableEnemy), 0, f-1)
	//create a stele if less than zhongli's max stele count and desired by player
	if (c.steleCount <= c.maxStele) && createStele {
		c.AddTask(func() {
			c.newStele(1860, max) //31 seconds
		}, "zhongli-create-stele", f)
	}

	//make a shield - enemy debuff arrows appear 3-5 frames after the damage number shows up in game
	c.AddTask(func() {
		c.addJadeShield()
	}, "zhongli-create-shield", f-1)
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)

	//deal damage when created
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Planet Befall",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Geo,
		Durability: 100,
		Mult:       burst[c.TalentLvlBurst()],
		FlatDmg:    0.33 * c.HPMax,
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, coretype.TargettableEnemy), f-1, f-1)

	if c.Base.Cons >= 2 {
		c.addJadeShield()
	}

	c.SetCDWithDelay(core.ActionBurst, 720, 6)
	c.ConsumeEnergy(6)
	return f, a
}
