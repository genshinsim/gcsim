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
		for i, mult := range shunsuiken[c.NormalCounter] {
			ai.Mult = mult[c.TalentLvlAttack()]
			c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), f-5+i, f-5+i)
			if c.Core.F > c.particleICD {
				c.particleICD = c.Core.F + 112 //best info we have rn
				c.QueueParticle("ayato", 1, core.Hydro, 80)
			}
		}
	} else {
		for i, mult := range attack[c.NormalCounter] {
			ai.Mult = mult[c.TalentLvlAttack()]
			c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), f-5+i, f-5+i)
		}

	}

	c.AdvanceNormalIndex()

	return f, a
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

	c.Core.Status.AddStatus("soukaikanka", 6*60+0) //doesn't account for animation
	c.Core.Log.NewEvent("Soukai Kanka acivated", core.LogCharacterEvent, c.Index, "expiry", c.Core.F+6*60+0)
	//figure out atk buff
	c.waterIllusion(ai, 6*60)
	c.SetCD(core.ActionSkill, 12*60)
	return f, a

}
func (c *char) waterIllusion(ai core.AttackInfo, delay int) {
	// currently assumes no attack
	c.AddTask(func() {
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, 0)
	}, "Water Illusion Burst", delay)
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
		// else {
		// 	c.ReduceActionCooldown(core.ActionSkill, 2*60)

		// 	return false
		// }
	}, "soukaiKankaProc")
}

func (c *char) Burst(p map[string]int) (int, int) {

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
	for delay := 0; delay < 12*60; delay += 30 {
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

	if c.Base.Cons >= 4 {
		val := make([]float64, core.EndStatType)
		val[core.DmgP] = 0.2
		for _, char := range c.Core.Chars {
			if char.CharIndex() == c.CharIndex() {
				continue
			}
			c.AddPreDamageMod(core.PreDamageMod{
				Key:    "ayato-c4",
				Expiry: 12 * 60,
				Amount: func(a *core.AttackEvent, t core.Target) ([]float64, bool) {
					if a.Info.AttackTag != core.AttackTagNormal {
						return nil, false
					}
					return val, true
				},
			})
		}
	}
	//add cooldown to sim
	c.SetCDWithDelay(core.ActionBurst, 15*60, 8)
	//use up energy
	c.ConsumeEnergy(8)

	return f, a
}

func (c *char) waveFlash() {
	val := make([]float64, core.EndStatType)
	val[core.DmgP] = skillpp[c.TalentLvlSkill()] * c.MaxHP()
	c.AddPreDamageMod(core.PreDamageMod{
		Key:    "ayato-waveFlash",
		Expiry: -1,
		Amount: func(a *core.AttackEvent, t core.Target) ([]float64, bool) {
			if a.Info.AttackTag != core.AttackTagNormal || c.Core.Status.Duration("soukaikanka") <= 0 {
				return nil, false
			}
			c.Core.Log.NewEvent("Waveflash Stacks: ", core.LogCharacterEvent, c.stacks, "expiry", c.Core.Status.Duration("soukaikanka"))
			return val, true
		},
	})
}
