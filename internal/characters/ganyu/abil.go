package ganyu

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
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
		AttackTag:  coretype.AttackTagNormal,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypePierce,
		Element:    core.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, coretype.TargettableEnemy), f, travel+f)

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
		bloom = 20
	}
	weakspot, ok := p["weakspot"]

	ai := core.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Frost Flake Arrow",
		AttackTag:    coretype.AttackTagExtra,
		ICDTag:       core.ICDTagNone,
		ICDGroup:     core.ICDGroupDefault,
		StrikeType:   core.StrikeTypePierce,
		Element:      coretype.Cryo,
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
		Element:    coretype.Cryo,
		Durability: 25,
		Mult:       lotus[c.TalentLvlSkill()],
	}

	snap := c.Snapshot(&ai)
	//flower damage immediately
	c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(2, false, coretype.TargettableEnemy), 30)
	//we get the orbs right away
	c.QueueParticle("ganyu", 2, coretype.Cryo, 90)

	//flower damage is after 6 seconds
	c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(2, false, coretype.TargettableEnemy), 360)
	// TODO: Particle flight time is 60s?
	c.QueueParticle("ganyu", 2, coretype.Cryo, 420)

	//add cooldown to sim
	// c.CD[charge] = c.Core.F + 10*60

	if c.Base.Cons == 6 {
		c.Core.AddStatus("ganyuc6", 1800)
	}

	c.SetCD(core.ActionSkill, 600)

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
		Element:    coretype.Cryo,
		Durability: 25,
		Mult:       shower[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)

	c.Core.AddStatus("ganyuburst", 15*60)

	rad, ok := p["radius"]
	if !ok {
		rad = 1
	}

	r := 2.5 + float64(rad)
	prob := r * r / 90.25

	lastHit := make(map[coretype.Target]int)
	// ccc := 0
	//tick every .3 sec, every fifth hit is targetted i.e. 1, 0, 0, 0, 0, 1
	for delay := 0; delay < 900; delay += 18 {
		c.AddTask(func() {
			//check if this hits first
			target := -1
			for i, t := range c.coretype.Targets {
				//skip for target 0 aka player
				if i == 0 {
					continue
				}
				if lastHit[t] < c.Core.Frame {
					target = i
					lastHit[t] = c.Core.Frame + 87 //cannot be targetted again for 1.45s
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
			c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(9, false, coretype.TargettableEnemy), 0)
		}, "ganyu-q", delay+f)

	}
	// c.AddTask(func() {
	// 	log.Println(ccc, prob)
	// }, "counts", 900+f+10)

	//a4 every .3 seconds for the duration of the burst, add ice dmg up to active char for 1sec
	//duration is 15 seconds
	//starts from end of cast
	for i := f; i < 900+f; i += 18 {
		t := i
		c.AddTask(func() {
			active := c.Core.Chars[c.Core.ActiveChar]
			val := make([]float64, core.EndStatType)
			val[coretype.CryoP] = 0.2
			active.AddMod(coretype.CharStatMod{
				Key: "ganyu-field",
				Amount: func() ([]float64, bool) {
					return val, true
				},
				Expiry: c.Core.Frame + 60,
			})
			if t >= 900-18 {
				c.coretype.Log.NewEvent("a4 last tick", coretype.LogCharacterEvent, c.Index, "ends_on", c.Core.Frame+60)
			}
		}, "ganyu-a4", i)
	}

	if c.Base.Cons >= 4 {
		//we just assume this lasts for the full duration since no one moves...
		start := c.Core.Frame

		val := make([]float64, core.EndStatType)
		c.AddMod(coretype.CharStatMod{
			Key:    "ganyu-c4",
			Expiry: c.Core.Frame + 1080,
			Amount: func() ([]float64, bool) {
				elapsed := c.Core.Frame - start
				stacks := int(elapsed / 180)
				if stacks > 5 {
					stacks = 5
				}
				val[core.DmgP] = float64(stacks) * 0.05
				return val, true
			},
		})
	}

	//add cooldown to sim
	c.SetCDWithDelay(core.ActionBurst, 15*60, 8)
	//use up energy
	c.ConsumeEnergy(8)

	return f, a
}
