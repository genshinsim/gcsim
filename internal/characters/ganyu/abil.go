package ganyu

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) Attack(p map[string]int) (int, int) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	f, a := c.ActionFrames(core.ActionAttack, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypePierce,
		Element:    core.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), f, travel+f)

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) Aimed(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAim, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	bloom, ok := p["bloom"]
	if !ok {
		bloom = 24
	}
	weakspot, ok := p["weakspot"]

	ai := core.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Frost Flake Arrow",
		AttackTag:    core.AttackTagExtra,
		ICDTag:       core.ICDTagNone,
		ICDGroup:     core.ICDGroupDefault,
		StrikeType:   core.StrikeTypePierce,
		Element:      core.Cryo,
		Durability:   25,
		Mult:         ffa[c.TalentLvlAttack()],
		HitWeakPoint: weakspot == 1,
	}

	// delay aim shot mostly to handle A1
	c.AddTask(func() {
		snap := c.Snapshot(&ai)
		if c.Core.F < c.a1Expiry {
			old := snap.Stats[core.CR]
			snap.Stats[core.CR] += .20
			c.Core.Log.NewEvent("a1 adding crit rate", core.LogCharacterEvent, c.Index, "old", old, "new", snap.Stats[core.CR], "expiry", c.a1Expiry)
		}

		c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefSingleTarget(1, core.TargettableEnemy), travel)

		ai.Abil = "Frost Flake Bloom"
		ai.Mult = ffb[c.TalentLvlAttack()]
		ai.HitWeakPoint = false
		c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(2, false, core.TargettableEnemy), travel+bloom)

		// first shot/bloom do not benefit from a1
		c.a1Expiry = c.Core.F + 60*5
	}, "ganyu-aim-snapshot", f)

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionSkill, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Ice Lotus",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Cryo,
		Durability: 25,
		Mult:       lotus[c.TalentLvlSkill()],
	}

	snap := c.Snapshot(&ai)
	//flower damage immediately
	c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(2, false, core.TargettableEnemy), 13)
	//we get the orbs right away
	c.QueueParticle("ganyu", 2, core.Cryo, 90)

	//flower damage is after 6 seconds
	c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(2, false, core.TargettableEnemy), 373)
	// TODO: Particle flight time is 60s?
	c.QueueParticle("ganyu", 2, core.Cryo, 420)

	//add cooldown to sim
	// c.CD[charge] = c.Core.F + 10*60

	if c.Base.Cons == 6 {
		c.Core.Status.AddStatus("ganyuc6", 1800)
	}

	c.SetCDWithDelay(core.ActionSkill, 600, 10)

	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionBurst, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Celestial Shower",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Cryo,
		Durability: 25,
		Mult:       shower[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)

	c.Core.Status.AddStatus("ganyuburst", 15*60+130)

	rad, ok := p["radius"]
	if !ok {
		rad = 1
	}

	r := 2.5 + float64(rad)
	prob := r * r / 90.25

	lastHit := make(map[core.Target]int)
	// ccc := 0
	//tick every .3 sec, every fifth hit is targetted i.e. 1, 0, 0, 0, 0, 1
	//first hit at 148
	for delay := a; delay < 900+a; delay += 18 {
		c.AddTask(func() {
			//check if this hits first
			target := -1
			for i, t := range c.Core.Targets {
				// skip non-enemy targets
				if t.Type() != core.TargettableEnemy {
					continue
				}
				if lastHit[t] < c.Core.F {
					target = i
					lastHit[t] = c.Core.F + 87 //cannot be targetted again for 1.45s
					break
				}
			}
			// log.Println(target)
			//[1:14 PM] Aluminum | Harbinger of Jank: assuming uniform distribution and enemy at center:
			//(radius_icicle + radius_enemy)^2 / radius_burst^2
			if target == -1 && c.Core.Rand.Float64() > prob {
				//no one getting hit
				return
			}
			//deal dmg
			c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(9, false, core.TargettableEnemy), 0)
		}, "ganyu-q", delay)

	}
	// c.AddTask(func() {
	// 	log.Println(ccc, prob)
	// }, "counts", 900+f+10)

	//a4 every .3 seconds for the duration of the burst, add ice dmg up to active char for 1sec
	//duration is 15 seconds
	//starts from end of cast
	mA4 := make([]float64, core.EndStatType)
	mC4 := make([]float64, core.EndStatType)
	mA4[core.CryoP] = 0.2
	for i := a; i < 900+a; i += 18 {
		t := i
		c.AddTask(func() {
			active := c.Core.Chars[c.Core.ActiveChar]
			active.AddMod(core.CharStatMod{
				Key:    "ganyu-field",
				Expiry: c.Core.F + 60,
				Amount: func() ([]float64, bool) {
					return mA4, true
				},
			})
			if t >= 900+a-18 {
				c.Core.Log.NewEvent("a4 last tick", core.LogCharacterEvent, c.Index, "ends_on", c.Core.F+60)
			}

			// C4: similar to A4 expect it lingers for 3s
			// assume this lasts for the full duration since no one moves...
			if c.Base.Cons >= 4 {
				// check for 1st tick and reset stacks if expired
				if t == a && !c.PreDamageModIsActive("ganyu-c4") {
					c.c4Stacks = 0
				}

				// increase stacks at 3s interval
				if (t-a)%180 == 0 {
					c.c4Stacks++
					if c.c4Stacks > 5 {
						c.c4Stacks = 5
					}
					mC4[core.DmgP] = float64(c.c4Stacks) * 0.05
				}

				// TODO: should be changed to target mod
				for _, char := range c.Core.Chars {
					char.AddPreDamageMod(core.PreDamageMod{
						Key:    "ganyu-c4",
						Expiry: c.Core.F + 60*3,
						Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
							return mC4, true
						},
					})
				}

				c.Core.Log.NewEvent("ganyu c4 stacks", core.LogCharacterEvent, c.Index, "stacks", c.c4Stacks)
			}
		}, "ganyu-burst-checks", i)
	}

	//add cooldown to sim
	c.SetCD(core.ActionBurst, 15*60)
	//use up energy
	c.ConsumeEnergy(3)

	return f, a
}
