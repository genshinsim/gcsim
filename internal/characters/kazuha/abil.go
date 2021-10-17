package kazuha

import "github.com/genshinsim/gsim/pkg/core"

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)

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
		c.QueueDmg(&d, f-2+i)
	}

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) HighPlungeAttack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionHighPlunge, p)
	ele := core.Physical
	if c.Core.LastAction.Target == "kazuha" && c.Core.LastAction.Typ == core.ActionSkill {
		ele = core.Anemo
	}

	_, ok := p["collide"]
	if ok {
		d := c.Snapshot(
			"Plunge",
			core.AttackTagPlunge,
			core.ICDTagNormalAttack,
			core.ICDGroupDefault,
			core.StrikeTypeSlash,
			ele,
			25,
			plunge[c.TalentLvlAttack()],
		)
		c.QueueDmg(&d, f-10)
	}

	//aoe dmg
	d := c.Snapshot(
		"Plunge",
		core.AttackTagPlunge,
		core.ICDTagNormalAttack,
		core.ICDGroupDefault,
		core.StrikeTypeSlash,
		ele,
		25,
		highPlunge[c.TalentLvlAttack()],
	)
	d.Targets = core.TargetAll
	c.QueueDmg(&d, f-8)

	// a2 if applies
	if c.a2Ele != core.NoElement {
		d := c.Snapshot(
			"Kazuha A2",
			core.AttackTagPlunge,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			c.a2Ele,
			25,
			2, //200%
		)
		d.Targets = core.TargetAll
		c.QueueDmg(&d, 10)
		c.a2Ele = core.NoElement
	}

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	hold := p["hold"]
	c.a2Ele = core.NoElement
	if hold == 0 {
		return c.skillPress(p)
	}
	return c.skillHold(p)
}

func (c *char) skillPress(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)
	d := c.Snapshot(
		"Chihayaburu",
		core.AttackTagElementalArt,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Anemo,
		25,
		skill[c.TalentLvlSkill()],
	)
	d.Targets = core.TargetAll
	c.QueueDmg(&d, f-10)

	c.QueueParticle("kazuha", 3, core.Anemo, 100)

	c.AddTask(c.absorbCheckA2(c.Core.F, 0, int(f/18)), "kaz-a2-absorb-check", 1)

	cd := 360
	if c.Base.Cons > 0 {
		cd = 324
	}
	if c.Base.Cons == 6 {
		c.c6Active = c.Core.F + f + 300
	}
	c.SetCD(core.ActionSkill, cd)

	return f, a
}

func (c *char) skillHold(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)
	d := c.Snapshot(
		"Chihayaburu",
		core.AttackTagElementalArt,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Anemo,
		50,
		skillHold[c.TalentLvlSkill()],
	)
	d.Targets = core.TargetAll
	c.QueueDmg(&d, f-10)

	c.QueueParticle("kazuha", 4, core.Anemo, 100)

	c.AddTask(c.absorbCheckA2(c.Core.F, 0, int(f/18)), "kaz-a2-absorb-check", 1)
	cd := 540
	if c.Base.Cons > 0 {
		cd = 486
	}
	if c.Base.Cons == 6 {
		c.c6Active = c.Core.F + f + 300
	}
	c.SetCD(core.ActionSkill, cd)
	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)

	c.qInfuse = core.NoElement

	d := c.Snapshot(
		"Kazuha Slash",
		core.AttackTagElementalBurst,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Anemo,
		50,
		burstSlash[c.TalentLvlBurst()],
	)
	d.Targets = core.TargetAll
	c.QueueDmg(&d, f-10)

	//apply dot and check for absorb
	d1 := d.Clone()
	d1.Abil = "Kazuha Slash (Dot)"
	d1.Mult = burstDot[c.TalentLvlBurst()]
	d1.Durability = 25

	c.AddTask(c.absorbCheckQ(c.Core.F, 0, int(310/18)), "kaz-absorb-check", 10)

	//424 start
	//493 first tick, 553, 612, 670, 729 <- so tick every second starting at 70 frames in
	for i := 70; i < 70+60*5; i += 60 {
		x := d1.Clone()
		c.AddTask(func() {
			if c.qInfuse != core.NoElement {
				d := c.Snapshot(
					"Kazuha Slash (Absorb Dot)",
					core.AttackTagElementalBurst,
					core.ICDTagNone,
					core.ICDGroupDefault,
					core.StrikeTypeDefault,
					c.qInfuse,
					25,
					burstEleDot[c.TalentLvlBurst()],
				)
				d.Targets = core.TargetAll
				c.Core.Combat.ApplyDamage(&d)
			}
			c.Core.Combat.ApplyDamage(&x)
		}, "kazuha-burst-tick", i)
	}

	//reset skill cd
	if c.Base.Cons > 0 {
		c.ResetActionCooldown(core.ActionSkill)
	}

	//add em to all char, but only activate if char is active
	if c.Base.Cons >= 2 {
		val := make([]float64, core.EndStatType)
		val[core.EM] = 200
		for _, char := range c.Core.Chars {
			this := char
			char.AddMod(core.CharStatMod{
				Key:    "kazuha-c2",
				Expiry: c.Core.F + 370,
				Amount: func(a core.AttackTag) ([]float64, bool) {
					if c.Core.ActiveChar != this.CharIndex() {
						return nil, false
					}
					return val, true
				},
			})
		}
	}

	if c.Base.Cons == 6 {
		c.c6Active = c.Core.F + f + 300
	}

	c.SetCD(core.ActionBurst, 15*60)
	c.Energy = 0
	return f, a
}

func (c *char) absorbCheckQ(src, count, max int) func() {
	return func() {
		if count == max {
			return
		}
		c.qInfuse = c.Core.AbsorbCheck(core.Pyro, core.Hydro, core.Electro, core.Cryo)

		// Special handling for Bennett field self-infusion while waiting for something comprehensive
		// Interaction is crucial to making many teams work correctly
		if c.Core.Status.Duration("btburst") > 0 {
			c.qInfuse = core.Pyro
		}

		if c.qInfuse != core.NoElement {
			return
		}
		//otherwise queue up
		c.AddTask(c.absorbCheckQ(src, count+1, max), "kaz-q-absorb-check", 18)
	}
}

func (c *char) absorbCheckA2(src, count, max int) func() {
	return func() {
		if count == max {
			return
		}
		c.a2Ele = c.Core.AbsorbCheck(core.Pyro, core.Hydro, core.Electro, core.Cryo)

		// Special handling for Bennett field self-infusion while waiting for something comprehensive
		// Interaction is crucial to making many teams work correctly
		if c.Core.Status.Duration("btburst") > 0 {
			c.a2Ele = core.Pyro
		}

		if c.a2Ele != core.NoElement {
			return
		}
		//otherwise queue up
		c.AddTask(c.absorbCheckA2(src, count+1, max), "kaz-a2-absorb-check", 18)
	}
}
