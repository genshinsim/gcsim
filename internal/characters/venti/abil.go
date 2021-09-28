package venti

import (
	"github.com/genshinsim/gsim/pkg/core"
)

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	for i, mult := range attack[c.NormalCounter] {
		d := c.Snapshot(
			//fmt.Sprintf("Normal %v", c.NormalCounter),
			"Normal",
			core.AttackTagNormal,
			core.ICDTagNormalAttack,
			core.ICDGroupDefault,
			core.StrikeTypeSlash,
			core.Physical,
			25,
			mult[c.TalentLvlAttack()],
		)
		c.QueueDmg(&d, f+travel+i)
	}

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) Aimed(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAim, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	d := c.Snapshot(
		"Aim (Charged)",
		core.AttackTagExtra,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypePierce,
		core.Anemo,
		25,
		aim[c.TalentLvlAttack()],
	)

	d.HitWeakPoint = true
	d.AnimationFrames = f

	c.QueueDmg(&d, travel+f)

	if c.Base.Cons >= 1 {
		d1 := d.Clone()
		d1.Mult = d.Mult / 3.0
		d1.Abil = "Aim (Charged) C1"
		d2 := d1.Clone()

		c.QueueDmg(&d1, travel+f)
		c.QueueDmg(&d2, travel+f)
	}

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionSkill, p)

	cd := 360
	d := c.Snapshot(
		"Skyward Sonnett",
		core.AttackTagElementalArt,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypePierce,
		core.Anemo,
		50,
		skillPress[c.TalentLvlSkill()],
	)
	d.Targets = core.TargetAll
	if p["hold"] == 1 {
		cd = 900
		d.Mult = skillHold[c.TalentLvlSkill()]
	}

	if c.Base.Cons >= 2 {
		c.applyC2(&d)
	}

	c.QueueDmg(&d, f-1)

	c.QueueParticle("venti", 4, core.Anemo, f+100)

	c.SetCD(core.ActionSkill, cd)
	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)

	c.qInfuse = core.NoElement

	//8 second duration, tick every .4 second
	d := c.Snapshot(
		"Wind's Grand Ode",
		core.AttackTagElementalBurst,
		core.ICDTagVentiBurstAnemo,
		core.ICDGroupVenti,
		core.StrikeTypeDefault,
		core.Anemo,
		25,
		burstDot[c.TalentLvlBurst()],
	)
	d.Targets = core.TargetAll

	if c.Base.Cons == 6 {
		c.applyC6(&d, core.Anemo)
	}

	for i := 24; i <= 480; i += 24 {
		x := d.Clone()
		c.QueueDmg(&x, i)
	}

	c.AddTask(c.absorbCheckQ(c.Core.F, 0, int(480/18)), "venti-absorb-check", 10)

	c.AddTask(func() {
		c.a4Restore()
	}, "venti-a4-restore", 480+f)

	c.SetCD(core.ActionBurst, 15*60)
	c.Energy = 0
	return f, a
}

func (c *char) a4Restore() {
	c.AddEnergy(15)

	if c.qInfuse != core.NoElement {
		for _, char := range c.Core.Chars {
			if char.Ele() == c.qInfuse {
				char.AddEnergy(15)
			}
		}
	}
}

func (c *char) burstInfusedTicks() {
	d := c.Snapshot(
		"Wind's Grand Ode (Infused)",
		core.AttackTagElementalBurst,
		core.ICDTagVentiBurstAnemo,
		core.ICDGroupVenti,
		core.StrikeTypeDefault,
		c.qInfuse,
		25,
		burstAbsorbDot[c.TalentLvlBurst()],
	)
	d.Targets = core.TargetAll

	if c.Base.Cons == 6 {
		c.applyC6(&d, c.qInfuse)
	}

	for i := 24; i <= 360; i += 24 {
		x := d.Clone()
		c.QueueDmg(&x, i)
	}
}

func (c *char) absorbCheckQ(src, count, max int) func() {
	return func() {
		if count == max {
			return
		}
		c.qInfuse = c.Core.AbsorbCheck(core.Pyro, core.Hydro, core.Electro, core.Cryo)
		if c.qInfuse != core.NoElement {
			//trigger dmg ticks here
			c.burstInfusedTicks()

			return
		}
		//otherwise queue up
		c.AddTask(c.absorbCheckQ(src, count+1, max), "venti-absorb-check", 18)
	}
}
