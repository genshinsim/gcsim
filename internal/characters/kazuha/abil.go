package kazuha

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Normal",
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeSlash,
		Element:    core.Physical,
		Durability: 25,
	}
	snap := c.Snapshot(&ai)
	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(0.3, false, core.TargettableEnemy), f-2+i)
	}

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) HighPlungeAttack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionHighPlunge, p)
	ele := core.Physical
	if c.Core.LastAction.Target == keys.Kazuha && c.Core.LastAction.Typ == core.ActionSkill {
		ele = core.Anemo
	}

	_, ok := p["collide"]
	if ok {
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Plunge",
			AttackTag:  core.AttackTagPlunge,
			ICDTag:     core.ICDTagNormalAttack,
			ICDGroup:   core.ICDGroupDefault,
			StrikeType: core.StrikeTypeSlash,
			Element:    ele,
			Durability: 25,
			Mult:       plunge[c.TalentLvlAttack()],
		}
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.3, false, core.TargettableEnemy), 0, f-10)
	}

	//aoe dmg
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Plunge",
		AttackTag:  core.AttackTagPlunge,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeSlash,
		Element:    ele,
		Durability: 25,
		Mult:       highPlunge[c.TalentLvlAttack()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1.5, false, core.TargettableEnemy), 0, f-8)

	// a2 if applies
	if c.a2Ele != core.NoElement {
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Kazuha A2",
			AttackTag:  core.AttackTagPlunge,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			StrikeType: core.StrikeTypeDefault,
			Element:    c.a2Ele,
			Durability: 25,
			Mult:       2,
		}

		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1.5, false, core.TargettableEnemy), 0, 10)
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
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Chihayaburu",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Anemo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1.5, false, core.TargettableEnemy), 0, f-10)

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
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Chihayaburu",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Anemo,
		Durability: 50,
		Mult:       skillHold[c.TalentLvlSkill()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1.5, false, core.TargettableEnemy), 0, f-10)

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
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Kazuha Slash",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Anemo,
		Durability: 50,
		Mult:       burstSlash[c.TalentLvlBurst()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1.5, false, core.TargettableEnemy), 0, f-10)

	//apply dot and check for absorb
	ai.Abil = "Kazuha Slash (Dot)"
	ai.Mult = burstDot[c.TalentLvlBurst()]
	ai.Durability = 25
	snap := c.Snapshot(&ai)

	c.AddTask(c.absorbCheckQ(c.Core.F, 0, int(310/18)), "kaz-absorb-check", 10)

	//424 start
	//493 first tick, 553, 612, 670, 729 <- so tick every second starting at 70 frames in
	for i := 70; i < 70+60*5; i += 60 {
		c.AddTask(func() {
			if c.qInfuse != core.NoElement {
				//TODO: does absorb dot tick snapshot?
				absorb := core.AttackInfo{
					ActorIndex: c.Index,
					Abil:       "Kazuha Slash (Absorb Dot)",
					AttackTag:  core.AttackTagElementalBurst,
					ICDTag:     core.ICDTagNone,
					ICDGroup:   core.ICDGroupDefault,
					StrikeType: core.StrikeTypeDefault,
					Element:    c.qInfuse,
					Durability: 25,
					Mult:       burstEleDot[c.TalentLvlBurst()],
				}
				c.Core.Combat.QueueAttack(absorb, core.NewDefCircHit(5, false, core.TargettableEnemy), 0, 0)
			}
			c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(5, false, core.TargettableEnemy), 0)
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
	c.ConsumeEnergy(7)
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
		// TODO: get rid of this once we add in self app
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
