package sucrose

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

var hitmarks = []int{17, 18, 28, 28}

func (c *char) Attack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAttack, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Anemo,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), 0, hitmarks[c.NormalCounter])

	c.AdvanceNormalIndex()

	if c.Base.Cons >= 4 {
		c.c4()
	}

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

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.3, false, core.TargettableEnemy), 0, 54)

	if c.Base.Cons >= 4 {
		c.c4()
	}

	return f, a
}

//sucrose's dash can be cancelled by her E and Q, so we override it here
func (c *char) Dash(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionDash, p)
	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)

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
		done = true
		c.a4()
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy, core.TargettableObject), 0, 42, cb)

	c.QueueParticle("sucrose", 4, core.Anemo, 150)

	//reduce charge by 1
	c.SetCDWithDelay(core.ActionSkill, eCD, 9)

	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)
	//tag a4
	//first hit at 137, then 113 frames between hits

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
		//lockout for 1 frame to prevent triggering multiple times on one attack
		if lockout > c.Core.F {
			return
		}
		lockout = c.Core.F + 1
		c.a4()
	}

	for i := 137; i <= duration+5; i += 113 {
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
	c.AddTask(c.absorbCheck(c.Core.F, 0, int(duration/18)), "absorb-check", 136)

	c.SetCDWithDelay(core.ActionBurst, 1200, 18)
	c.ConsumeEnergy(21)
	return f, a
}

func (c *char) absorbCheck(src, count, max int) func() {
	return func() {
		if count == max {
			return
		}
		c.qInfused = c.Core.AbsorbCheck(c.infuseCheckLocation, core.Pyro, core.Hydro, core.Electro, core.Cryo)

		if c.qInfused != core.NoElement {
			if c.Base.Cons >= 6 {
				c.c6()
			}
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
