package venti

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Normal",
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
	}

	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 0, f+travel+i)
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

	ai := core.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Aim (Charged)",
		AttackTag:    core.AttackTagExtra,
		ICDTag:       core.ICDTagNone,
		ICDGroup:     core.ICDGroupDefault,
		Element:      core.Anemo,
		Durability:   25,
		Mult:         aim[c.TalentLvlAttack()],
		HitWeakPoint: true,
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
		ai.Mult = skillHold[c.TalentLvlSkill()]
	}

	var cb core.AttackCBFunc

	if c.Base.Cons >= 2 {
		cb = c2cb
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(4, false, core.TargettableEnemy), 0, f-1, cb)

	c.QueueParticle("venti", 4, core.Anemo, f+100)

	c.SetCD(core.ActionSkill, cd)
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
	snap := c.Snapshot(&ai)

	var cb core.AttackCBFunc
	if c.Base.Cons == 6 {
		cb = c6cb(core.Anemo)
	}

	for i := 24; i <= 480; i += 24 {
		c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(4, false, core.TargettableEnemy), i, cb)
	}

	c.AddTask(c.absorbCheckQ(c.Core.F, 0, int(480/18)), "venti-absorb-check", 10)

	c.AddTask(func() {
		c.a4Restore()
	}, "venti-a4-restore", 480+f)

	c.SetCD(core.ActionBurst, 15*60)
	c.ConsumeEnergy(90)
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
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Wind's Grand Ode (Infused)",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagVentiBurstAnemo,
		ICDGroup:   core.ICDGroupVenti,
		Element:    c.qInfuse,
		Durability: 25,
		Mult:       burstAbsorbDot[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)
	var cb core.AttackCBFunc
	if c.Base.Cons == 6 {
		cb = c6cb(c.qInfuse)
	}

	for i := 24; i <= 360; i += 24 {
		c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(4, false, core.TargettableEnemy), i, cb)
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
