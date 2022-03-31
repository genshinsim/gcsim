package ayato

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

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
		ai.Mult = shunsuiken[c.NormalCounter][c.TalentLvlAttack()]
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), f, f, c.generateParticles)
	} else {
		for i, mult := range attack[c.NormalCounter] {
			ai.Mult = mult[c.TalentLvlAttack()]
			c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), f-5+i, f-5+i)
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

	for i := 0; i < 3; i++ {
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), f-3+i, f-3+i)
	}

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	delay := p["delay"]
	if delay < 30 { //this might be too low? might be 48?
		delay = 30
	}
	if delay > 6*60 {
		delay = 360
	}
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
	}, "Water Illusion Burst", delay)

	c.Core.Status.AddStatus("soukaikanka", 6*60+0) //doesn't account for animation
	c.Core.Log.NewEvent("Soukai Kanka acivated", core.LogCharacterEvent, c.Index, "expiry", c.Core.F+6*60+0)
	//figure out atk buff
	if c.Base.Cons >= 6 {
		c.c6ready = true

	}
	c.SetCD(core.ActionSkill, 12*60)
	return f, a

}

func (c *char) soukaiKankaHook() {
	c.Core.Events.Subscribe(core.EventType(core.OnDamage), func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)

		if c.Core.Status.Duration("soukaikanka") <= 0 {
			return false
		}
		if atk.Info.AttackTag != core.AttackTagNormal {
			return false
		}

		if atk.Info.ActorIndex == c.Index {
			c.stacks++
			c.Core.Log.NewEvent("Soukai Kanka Proc'd by", core.LogCharacterEvent, c.Index)
			if c.stacks > c.stacksMax {
				c.stacks = c.stacksMax
			}
			return false
		}
		return false
	}, "soukaiKankaProc")
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
	snap := c.Snapshot(&ai)

	rad, ok := p["radius"]
	if !ok {
		rad = 1
	}

	r := 2.5 + float64(rad)
	prob := r * r / 90.25

	lastHit := make(map[core.Target]int)
	// ccc := 0
	//tick every .3 sec, every fifth hit is targetted i.e. 1, 0, 0, 0, 0, 1
	for delay := 0; delay < dur*60; delay += 30 {
		c.AddTask(func() {
			//check if this hits first
			target := -1
			for i, t := range c.Core.Targets {
				//skip for target 0 aka player
				if i == 0 {
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
		}, "ayato-q", delay+f)

	}

	c.Core.Status.AddStatus("ayatoburst", dur*60) //doesn't account for animation
	// if c.Base.Cons >= 4 {
	val := make([]float64, core.EndStatType)
	val[core.DmgP] = 0.2
	for _, char := range c.Core.Chars {
		if char.CharIndex() == c.CharIndex() {
			continue
		}
		c.AddPreDamageMod(core.PreDamageMod{
			Key:    "ayato-c4",
			Expiry: dur * 60,
			Amount: func(a *core.AttackEvent, t core.Target) ([]float64, bool) {
				if a.Info.AttackTag != core.AttackTagNormal {
					return nil, false
				}
				return val, true
			},
		})
	}
	// }
	//add cooldown to sim
	c.SetCDWithDelay(core.ActionBurst, 20*60, 8)
	//use up energy
	c.ConsumeEnergy(8)

	return f, a
}

func (c *char) namisenStack() {
	c.Core.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		ae := args[1].(*core.AttackEvent)

		if c.Core.ActiveChar != c.CharIndex() {
			return false
		}
		if ae.Info.AttackTag != core.AttackTagNormal {
			return false
		}
		if c.Core.Status.Duration("soukaikanka") <= 0 {
			return false
		}

		stats, _ := c.SnapshotStats()
		dmgAdded := (c.Base.HP*(1+stats[core.HPP]) + stats[core.HP]) * skillpp[c.TalentLvlSkill()] * float64(c.stacks)
		ae.Info.FlatDmg += dmgAdded

		c.Core.Log.NewEvent("Waveflash Stacks: ", core.LogCharacterEvent, c.stacks, "expiry", c.Core.Status.Duration("soukaikanka"))
		return false
	}, "ayato-waveflash")
}
