package zhongli

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

func (c *char) Attack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAttack, p)
	d := c.Snapshot(
		fmt.Sprintf("Normal %v", c.NormalCounter),
		core.AttackTagNormal,
		core.ICDTagNormalAttack,
		core.ICDGroupDefault,
		core.StrikeTypeSpear,
		core.Physical,
		25,
		attack[c.NormalCounter][c.TalentLvlAttack()],
	)
	d.FlatDmg = 0.0139 * c.HPMax

	for i := 0; i < hits[c.NormalCounter]; i++ {
		x := d.Clone()
		c.QueueDmg(&x, f-i)
	}

	c.AdvanceNormalIndex()
	return f, a
}

func (c *char) ChargeAttack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionCharge, p)
	d := c.Snapshot(
		"Charge",
		core.AttackTagExtra,
		core.ICDTagExtraAttack,
		core.ICDGroupPole,
		core.StrikeTypeSpear,
		core.Physical,
		25,
		charge[c.TalentLvlAttack()],
	)
	d.FlatDmg = 0.0139 * c.HPMax

	c.QueueDmg(&d, f-1)

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {

	cd := 240
	f, a := c.ActionFrames(core.ActionSkill, p)

	max, ok := p["max"]
	if !ok {
		max = 3
	}

	//press does no dmg
	if p["hold"] == 1 {
		c.skillHold(f, max, true)
		cd = 720
	} else if p["hold_nosteele"] == 1 {
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
	c.newSteele(f, 1860, max)
}

func (c *char) skillHold(f, max int, createSteele bool) {
	//hold does dmg
	d := c.Snapshot(
		"Stone Stele (Hold)",
		core.AttackTagElementalArt,
		core.ICDTagElementalArt,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Geo,
		25,
		skillHold[c.TalentLvlSkill()],
	)
	d.FlatDmg = 0.019 * c.HPMax
	d.Targets = core.TargetAll

	c.QueueDmg(&d, f-1)

	//create a steele if none exists and desired by player
	if (c.steeleCount == 0) && createSteele {
		c.newSteele(f, 1860, max) //31 seconds
	}

	//make a shield
	c.addJadeShield()
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)

	//deal damage when created
	d := c.Snapshot(
		"Planet Befall",
		core.AttackTagElementalBurst,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Geo,
		100,
		burst[c.TalentLvlBurst()],
	)
	d.Targets = core.TargetAll
	d.FlatDmg = 0.33 * c.HPMax

	c.QueueDmg(&d, f-1)

	if c.Base.Cons >= 2 {
		c.addJadeShield()
	}

	c.SetCD(core.ActionBurst, 720)
	c.Energy = 0
	return f, a
}
