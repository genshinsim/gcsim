package sucrose

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) Attack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAttack, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Normal",
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Anemo,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), 0, f-1)

	c.AdvanceNormalIndex()

	c.c4()

	return f, a
}

func (c *char) ChargeAttack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionCharge, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Anemo,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.3, false, core.TargettableEnemy), 0, f-1)

	if c.Base.Cons >= 4 {
		c.c4()
	}

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)
	//41 frame delay
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Astable Anemohypostasis Creation-6308",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Anemo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}

	done := false
	cb := func(a core.AttackCB) {
		if done {
			return
		}
		c.Core.Status.AddStatus("sucrosea4", 480)
		c.a4EM[core.EM] = 0.2 * c.Stat(core.EM)
		c.Core.Log.Debugw("sucrose a4 triggered", "frame", c.Core.F, "event", core.LogCharacterEvent, "em snapshot", c.a4EM, "expiry", c.Core.F+480)
		done = true
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), 0, 41, cb)

	c.QueueParticle("sucrose", 4, core.Anemo, 150)

	//reduce charge by 1
	c.SetCD(core.ActionSkill, eCD)

	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)
	//tag a4
	//3 hits, 135, 249, 368; all 3 applied swirl; c2 i guess adds 2 second so one more hit
	//let's just assume 120, 240, 360, 480

	duration := 360
	if c.Base.Cons >= 2 {
		duration = 480
	}

	c.qInfused = core.NoElement

	// c.S.Status["sucroseburst"] = c.Core.F + count
	c.Core.Status.AddStatus("sucroseburst", duration)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Forbidden Creation-Isomer 75/Type II",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Anemo,
		Durability: 25,
		Mult:       burstDot[c.TalentLvlBurst()],
	}
	//TODO: does sucrose burst snapshot?
	snap := c.Snapshot(&ai)
	//TODO: does burst absorb snapshot
	aiAbs := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Forbidden Creation-Isomer 75/Type II (Absorb)",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.NoElement,
		Durability: 25,
		Mult:       burstAbsorb[c.TalentLvlBurst()],
	}
	snapAbs := c.Snapshot(&aiAbs)

	lockout := 0
	cb := func(a core.AttackCB) {
		if lockout > c.Core.F {
			return
		}
		c.Core.Status.AddStatus("sucrosea4", 480)
		c.a4EM[core.EM] = 0.2 * c.Stat(core.EM)
		c.Core.Log.Debugw("sucrose a4 triggered", "frame", c.Core.F, "event", core.LogCharacterEvent, "em snapshot", c.a4EM, "expiry", c.Core.F+480)
		//lockout for 1 frame to prevent triggering multiple times on one attack
		lockout = c.Core.F + 1
	}

	for i := 120; i <= duration; i += 120 {
		c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(5, false, core.TargettableEnemy), i, cb)

		c.AddTask(func() {
			if c.qInfused != core.NoElement {
				aiAbs.Element = c.qInfused
				c.Core.Combat.QueueAttackWithSnap(aiAbs, snapAbs, core.NewDefCircHit(5, false, core.TargettableEnemy), 0)
			}
			//check if infused
		}, "sucrose-burst-em", i)
	}

	//
	c.AddTask(c.absorbCheck(c.Core.F, 0, int(duration/18)), "absorb-check", f)

	c.SetCD(core.ActionBurst, 1200)
	c.ConsumeEnergy(26)
	return f, a
}

func (c *char) absorbCheck(src, count, max int) func() {
	return func() {
		if count == max {
			return
		}
		c.qInfused = c.Core.AbsorbCheck(core.Pyro, core.Hydro, core.Electro, core.Cryo)

		if c.qInfused != core.NoElement {
			return
		}
		//otherwise queue up
		c.AddTask(c.absorbCheck(src, count+1, max), "sucrose-absorb-check", 18)
	}
}

// func (c *char) absorbCheck(src int, count int, max int) func() {
// 	return func() {
// 		//max number of scans reached
// 		if count == max {
// 			return
// 		}

// 		fire := false
// 		water := false
// 		electric := false
// 		ice := false

// 		//scan through all targets, order is fire > water > electric > ice/frozen
// 		for _, t := range c.Core.Targets {
// 			switch t.AuraType() {
// 			case core.Pyro:
// 				fire = true
// 			case core.Hydro:
// 				water = true
// 			case core.Electro:
// 				electric = true
// 			case core.Cryo:
// 				ice = true
// 			case core.EC:
// 				water = true
// 			case core.Frozen:
// 				ice = true
// 			}
// 		}

// 		switch {
// 		case fire:
// 			c.qInfused = core.Pyro
// 		case water:
// 			c.qInfused = core.Hydro
// 		case electric:
// 			c.qInfused = core.Electro
// 		case ice:
// 			c.qInfused = core.Cryo
// 		default:
// 			//nothing found, queue next
// 			c.AddTask(c.absorbCheck(src, count+1, max), "absorb-detect", 18) //every 0.3 seconds
// 		}
// 	}
// }
