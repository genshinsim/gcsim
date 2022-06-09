package klee

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

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
		ICDTag:     core.ICDTagKleeFireDamage,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Pyro,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	cb := func(a core.AttackCB) { c.a1() }
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), f, f+travel, cb)

	c.c1(f + travel)

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) ChargeAttack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionCharge, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge",
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Pyro,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}
	snap := c.Snapshot(&ai)

	//stam is calculated before this func is called so it's safe to
	//set spark to 0 here
	if c.Core.Status.Duration("kleespark") > 0 {
		snap.Stats[core.DmgP] += .50
		c.Core.Status.DeleteStatus("kleespark")
		c.Core.Log.NewEvent("klee consumed spark", core.LogCharacterEvent, c.Index, "icd", c.sparkICD)
	}

	c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(2, false, core.TargettableEnemy), f+travel)

	c.c1(f + travel)

	return f, a
}

// Has two parameters, "bounce" determines the number of bounces that hit
// "mine" determines the number of mines that hit the enemy
func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)

	bounce, ok := p["bounce"]
	if !ok {
		bounce = 1
	}

	//mine lives for 5 seconds
	//3 bounces, roughly 30, 70, 110 hits
	cb := func(a core.AttackCB) { c.a1() }
	for i := 0; i < bounce; i++ {
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Jumpy Dumpty",
			AttackTag:  core.AttackTagElementalArt,
			ICDTag:     core.ICDTagKleeFireDamage,
			ICDGroup:   core.ICDGroupDefault,
			StrikeType: core.StrikeTypeBlunt,
			Element:    core.Pyro,
			Durability: 25,
			Mult:       jumpy[c.TalentLvlSkill()],
		}

		// 3rd bounce is 2B
		if i == 2 {
			ai.Durability = 50
		}

		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, f+30+i*40, cb)
	}

	if bounce > 0 {
		c.QueueParticle("klee", 4, core.Pyro, 130)
	}

	minehits, ok := p["mine"]
	if !ok {
		minehits = 2
	}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Jumpy Dumpty Mine Hit",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagKleeFireDamage,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Pyro,
		Durability: 25,
		Mult:       mine[c.TalentLvlSkill()],
	}

	var c2cb func(a core.AttackCB)

	if c.Base.Cons >= 2 {
		c2cb = func(a core.AttackCB) {
			a.Target.AddDefMod("kleec2", -.233, 600)
		}
	}

	//roughly 160 frames after mines are laid
	for i := 0; i < minehits; i++ {
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), 0, f+160, c2cb)
	}

	c.c1(f + 30)

	c.SetCD(core.ActionSkill, 1200)

	// c.SetCD(def.ActionSkill, 20*60)
	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Sparks'n'Splash",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Pyro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
		NoImpulse:  true,
	}
	//lasts 10 seconds, starts after 2.2 seconds maybe?
	c.Core.Status.AddStatus("kleeq", 600+132)

	//every 1.8 second +on added shoots between 3 to 5, ignore the queue thing.. space it out .2 between each wave i guess

	// snapshot at end of animation?
	var snap core.Snapshot
	c.AddTask(func() {
		snap = c.Snapshot(&ai)
	}, "klee-burst-snapshot", 100)

	for i := 132; i < 732; i += 108 {
		c.AddTask(func() {
			//no more if burst has ended early
			if c.Core.Status.Duration("kleeq") <= 0 {
				return
			}
			//wave 1 = 1
			c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(1, false, core.TargettableEnemy), 0)
			//wave 2 = 1 + 30% chance of 1
			c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(1, false, core.TargettableEnemy), 12)
			if c.Core.Rand.Float64() < 0.3 {
				c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(1, false, core.TargettableEnemy), 12)
			}
			//wave 3 = 1 + 50% chance of 1
			c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(1, false, core.TargettableEnemy), 24)
			if c.Core.Rand.Float64() < 0.5 {
				c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(1, false, core.TargettableEnemy), 24)
			}
		}, "klee-burst", i)
	}

	//every 3 seconds add energy if c6
	if c.Base.Cons == 6 {
		for i := f + 180; i < f+600; i += 180 {
			c.AddTask(func() {
				//no more if burst has ended early
				if c.Core.Status.Duration("kleeq") <= 0 {
					return
				}

				for i, x := range c.Core.Chars {
					if i == c.Index {
						continue
					}
					x.AddEnergy("klee-c6", 3)
				}
			}, "klee-c6", i)
		}

		// add 10% pyro for 25s
		m := make([]float64, core.EndStatType)
		m[core.PyroP] = .1
		for _, x := range c.Core.Chars {
			x.AddMod(core.CharStatMod{
				Key:    "klee-c6",
				Amount: func() ([]float64, bool) { return m, true },
				Expiry: c.Core.F + 1500,
			})
		}
	}

	c.c1(132)

	c.SetCDWithDelay(core.ActionBurst, 15*60, 15)
	c.ConsumeEnergy(15)
	return f, a
}
