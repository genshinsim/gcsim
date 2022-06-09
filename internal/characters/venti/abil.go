package venti

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

var hitmarks = [][]int{{17, 27}, {19}, {28}, {15, 28}, {17}, {49}}

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
	}

	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		// TODO - double check snapshotDelay
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), hitmarks[c.NormalCounter][i], hitmarks[c.NormalCounter][i]+travel)
	}

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) Aimed(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAim, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot, ok := p["weakspot"]

	ai := core.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Aim (Charged)",
		AttackTag:    core.AttackTagExtra,
		ICDTag:       core.ICDTagNone,
		ICDGroup:     core.ICDGroupDefault,
		Element:      core.Anemo,
		Durability:   25,
		Mult:         aim[c.TalentLvlAttack()],
		HitWeakPoint: weakspot == 1,
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(.1, false, core.TargettableEnemy), f, travel+f)

	if c.Base.Cons >= 1 {
		ai.Abil = "Aim (Charged) C1"
		ai.Mult = ai.Mult / 3.0
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(.1, false, core.TargettableEnemy), f, travel+f)
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(.1, false, core.TargettableEnemy), f, travel+f)
	}

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionSkill, p)

	cd := 360
	cdstart := 21
	hitmark := 51
	ai := core.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Skyward Sonnett",
		AttackTag:    core.AttackTagElementalArt,
		ICDTag:       core.ICDTagNone,
		ICDGroup:     core.ICDGroupDefault,
		Element:      core.Anemo,
		Durability:   50,
		Mult:         skillPress[c.TalentLvlSkill()],
		HitWeakPoint: true,
	}

	if p["hold"] == 1 {
		cd = 900
		cdstart = 34
		hitmark = 74
		ai.Mult = skillHold[c.TalentLvlSkill()]
	}

	var cb core.AttackCBFunc

	if c.Base.Cons >= 2 {
		cb = c2cb
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(4, false, core.TargettableEnemy), 0, hitmark, cb)

	c.QueueParticle("venti", 3, core.Anemo, hitmark+100)

	c.SetCDWithDelay(core.ActionSkill, cd, cdstart)
	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)

	c.qInfuse = core.NoElement

	//8 second duration, tick every .4 second
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Wind's Grand Ode",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagVentiBurstAnemo,
		ICDGroup:   core.ICDGroupVenti,
		Element:    core.Anemo,
		Durability: 25,
		Mult:       burstDot[c.TalentLvlBurst()],
	}
	c.aiAbsorb = ai
	c.aiAbsorb.Abil = "Wind's Grand Ode (Infused)"
	c.aiAbsorb.Mult = burstAbsorbDot[c.TalentLvlBurst()]
	c.aiAbsorb.Element = core.NoElement

	// snapshot is around cd frame and 1st tick?
	var snap core.Snapshot
	c.AddTask(func() {
		snap = c.Snapshot(&ai)
		c.snapAbsorb = c.Snapshot(&c.aiAbsorb)
	}, "venti-q-snapshot", 104)

	var cb core.AttackCBFunc
	if c.Base.Cons >= 6 {
		cb = c6cb(core.Anemo)
	}

	// starts at 106 with 24f interval between ticks. 20 total
	for i := 0; i < 20; i++ {
		c.AddTask(func() {
			c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(4, false, core.TargettableEnemy), 0, cb)
		}, "venti-burst-tick", 106+24*i)
	}
	// Infusion usually occurs after 4 ticks of anemo according to KQM library
	c.AddTask(c.absorbCheckQ(c.Core.F, 0, int((480-24*4)/18)), "venti-absorb-check", 106+24*3)

	c.AddTask(func() {
		c.a4Restore()
	}, "venti-a4-restore", 480+f)

	c.SetCDWithDelay(core.ActionBurst, 15*60, 81)
	c.ConsumeEnergy(84)
	return f, a
}

func (c *char) HighPlungeAttack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionHighPlunge, p)

	ai := core.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Plunge",
		AttackTag:      core.AttackTagPlunge,
		ICDTag:         core.ICDTagNormalAttack,
		ICDGroup:       core.ICDGroupDefault,
		StrikeType:     core.StrikeTypeBlunt,
		Element:        core.Physical,
		Durability:     25,
		Mult:           highPlunge[c.TalentLvlAttack()],
		IgnoreInfusion: true,
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1.5, false, core.TargettableEnemy), f, f)

	return f, a
}

func (c *char) a4Restore() {
	c.AddEnergy("venti-a4", 15)

	if c.qInfuse != core.NoElement {
		for _, char := range c.Core.Chars {
			if char.Ele() == c.qInfuse {
				char.AddEnergy("venti-a4", 15)
			}
		}
	}
}

func (c *char) burstInfusedTicks() {
	var cb core.AttackCBFunc
	if c.Base.Cons >= 6 {
		cb = c6cb(c.qInfuse)
	}

	// ticks at 24f. 15 total
	for i := 0; i < 15; i++ {
		c.Core.Combat.QueueAttackWithSnap(c.aiAbsorb, c.snapAbsorb, core.NewDefCircHit(4, false, core.TargettableEnemy), i*24, cb)
	}
}

func (c *char) absorbCheckQ(src, count, max int) func() {
	return func() {
		if count == max {
			return
		}
		c.qInfuse = c.Core.AbsorbCheck(c.infuseCheckLocation, core.Pyro, core.Hydro, core.Electro, core.Cryo)
		if c.qInfuse != core.NoElement {
			c.aiAbsorb.Element = c.qInfuse
			switch c.qInfuse {
			case core.Pyro:
				c.aiAbsorb.ICDTag = core.ICDTagVentiBurstPyro
			case core.Hydro:
				c.aiAbsorb.ICDTag = core.ICDTagVentiBurstHydro
			case core.Electro:
				c.aiAbsorb.ICDTag = core.ICDTagVentiBurstElectro
			case core.Cryo:
				c.aiAbsorb.ICDTag = core.ICDTagVentiBurstCryo
			}
			//trigger dmg ticks here
			c.burstInfusedTicks()
			return
		}
		//otherwise queue up
		c.AddTask(c.absorbCheckQ(src, count+1, max), "venti-absorb-check", 18)
	}
}
