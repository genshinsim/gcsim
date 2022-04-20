package ayato

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

var hitmarks = [][]int{{12}, {18}, {20}, {22, 25}, {41}}

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)

	ai := core.AttackInfo{
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
	}
	if c.Core.Status.Duration("soukaikanka") > 0 {
		ai.Mult = shunsuiken[c.NormalCounter][c.TalentLvlSkill()]
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, f, c.generateParticles, c.skillStacks)
	} else {
		for i, mult := range attack[c.NormalCounter] {
			ai.Mult = mult[c.TalentLvlAttack()]
			c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), hitmarks[c.NormalCounter][i], hitmarks[c.NormalCounter][i])
		}
	}

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) generateParticles(ac core.AttackCB) {
	if c.Core.F > c.particleICD {
		c.particleICD = c.Core.F + 114
		count := 1
		if c.Core.Rand.Float64() < 0.5 {
			count++
		}
		c.QueueParticle("ayato", count, core.Hydro, 80)
	}
}

func (c *char) skillStacks(ac core.AttackCB) {
	if c.stacks < c.stacksMax {
		c.stacks++
		c.Core.Log.NewEvent("gained namisen stack", core.LogCharacterEvent, c.Index, "stacks", c.stacks)
	}
}

func (c *char) ChargeAttack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionCharge, p)

	ai := core.AttackInfo{
		Abil:       "Charge",
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagExtraAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
		Mult:       ca[c.TalentLvlAttack()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), f, f)

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	delay := p["illusion_delay"]
	if delay < 35 {
		delay = 35
	}
	if delay > 6*60 {
		delay = 360
	}
	hitlag := p["hitlag_extend"]
	f, a := c.ActionFrames(core.ActionSkill, p)
	ai := core.AttackInfo{
		Abil:       "Kamisato Art: Kyouka",
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagElementalArt,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Hydro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	c.AddTask(func() {
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(3.5, false, core.TargettableEnemy), 0, 0)
		//add a namisen stack
		if c.stacks < c.stacksMax {
			c.stacks++
		}
	}, "Water Illusion Burst", delay)

	c.Core.Status.AddStatus("soukaikanka", 6*60+f+hitlag) //add animation to the duration
	c.Core.Log.NewEvent("Soukai Kanka acivated", core.LogCharacterEvent, c.Index, "expiry", c.Core.F+6*60+f+hitlag)
	//figure out atk buff
	if c.Base.Cons >= 6 {
		c.c6ready = true

	}
	c.SetCD(core.ActionSkill, 12*60)
	return f, a

}

func (c *char) Burst(p map[string]int) (int, int) {

	dur := 18
	f, a := c.ActionFrames(core.ActionBurst, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Kamisato Art: Suiyuu",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Hydro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	// snapshot when the circle forms (is this correct?)
	var snap core.Snapshot
	c.AddTask(func() { snap = c.Snapshot(&ai) }, "ayato-q-snapshot", f)

	rad, ok := p["radius"]
	if !ok {
		rad = 1
	}

	r := 2.5 + float64(rad)
	prob := r * r / 90.25

	lastHit := make(map[core.Target]int)
	// ccc := 0
	//tick every .5 sec, every fourth hit is targetted i.e. 1, 0, 0, 0, 1
	for delay := 0; delay < dur*60; delay += 30 {
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
					lastHit[t] = c.Core.F + 117 //cannot be targetted again for 1.95s
					break
				}
			}
			// log.Println(target)
			//[1:14 PM] Aluminum | Harbinger of Jank: assuming uniform distribution and enemy at center:
			//(radius_droplet + radius_enemy)^2 / radius_burst^2
			if target == -1 && c.Core.Rand.Float64() > prob {
				//no one getting hit
				return
			}
			//deal dmg
			c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(9, false, core.TargettableEnemy), 0)
		}, "ayato-q", delay+139)

	}

	c.Core.Status.AddStatus("ayatoburst", dur*60+f)

	// NA buff starts after cast, ticks every 0.5s and last for 1.5s
	m := make([]float64, core.EndStatType)
	m[core.DmgP] = burstatkp[c.TalentLvlBurst()]
	for i := f; i < f+dur*60; i += 30 {
		c.AddTask(func() {
			active := c.Core.Chars[c.Core.ActiveChar]
			active.AddPreDamageMod(core.PreDamageMod{
				Key:    "ayato-burst",
				Expiry: c.Core.F + 90,
				Amount: func(a *core.AttackEvent, t core.Target) ([]float64, bool) {
					return m, a.Info.AttackTag == core.AttackTagNormal
				},
			})
		}, "ayato-burst-buff", i)
	}

	if c.Base.Cons >= 4 {
		m := make([]float64, core.EndStatType)
		m[core.AtkSpd] = 0.15
		for _, char := range c.Core.Chars {
			char.AddMod(core.CharStatMod{
				Key:    "ayato-c4",
				Expiry: c.Core.F + 15*60,
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})
		}
	}
	//add cooldown to sim
	c.SetCD(core.ActionBurst, 20*60)
	//use up energy
	c.ConsumeEnergy(5)

	return f, a
}
