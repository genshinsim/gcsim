package heizou

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

//var hitmarks = []int{17, 18, 28, 28}

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
	}
	for _, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		// TODO - double check snapshotDelay
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f, f)
	}

	c.AdvanceNormalIndex()

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

	return f, a
}

// //heizou's dash can be cancelled by her E and Q, so we override it here
// func (c *char) Dash(p map[string]int) (int, int) {
// 	f, a := c.ActionFrames(core.ActionDash, p)
// 	return f, a
// }

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)
	stackDelay := 45 //delay per stack while holding
	//TODO: Verify attack frame
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Heartstopper Strike",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Anemo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()] + float64(c.decStack)*decBonus[c.TalentLvlSkill()],
	}
	//a4 trigger
	done := false
	cb := func(a core.AttackCB) {
		if done {
			return
		}
		done = true
		switch c.decStack {
		case 2, 3:
			count := 2
			if c.Core.Rand.Float64() < .5 {
				count++
			}
			c.QueueParticle("heizou", count, core.Anemo, 150) //50% to generate 3 particles

		case 4:
			c.QueueParticle("heizou", 3, core.Anemo, 150)
		default:
			c.QueueParticle("heizou", 2, core.Anemo, 150) //2 particles on 0-1 stacks
		}
		c.AddTask(func() {
			c.decStack = 0
			c.Core.Log.NewEvent(
				"stack removed",
				core.LogCharacterEvent,
				c.Index,
			)
		}, "remove stack", 1)
		c.a4()
	}
	stackToBeAdded := p["hold"] //guarantee that you dont surpass more than 4 stacks with holding
	if stackToBeAdded+c.decStack > 4 {
		diff := 4 - c.decStack
		stackToBeAdded = diff
	}

	switch stackToBeAdded {
	case 1, 2, 3, 4:
		c.decStack += stackToBeAdded
		if c.Base.Cons >= 6 {
			c.c6()
		}
		if c.decStack > 3 {
			f = 44 + stackDelay*stackToBeAdded
			a = 44 + stackDelay*stackToBeAdded
			ai.Mult = skill[c.TalentLvlSkill()] + float64(c.decStack)*decBonus[c.TalentLvlSkill()] + convicBonus[c.TalentLvlSkill()]
			c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy, core.TargettableObject), f, a, cb)
		} else {
			f = 28 + stackDelay*stackToBeAdded
			a = 28 + stackDelay*stackToBeAdded
			ai.Mult = skill[c.TalentLvlSkill()] + float64(c.decStack)*decBonus[c.TalentLvlSkill()]
			c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy, core.TargettableObject), f, a, cb)
		}

	default:
		if c.Base.Cons >= 6 {
			c.c6()
		}
		if c.decStack > 3 {
			ai.Mult = skill[c.TalentLvlSkill()] + float64(c.decStack)*decBonus[c.TalentLvlSkill()] + convicBonus[c.TalentLvlSkill()]
		}

		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy, core.TargettableObject), f, f, cb)
	}
	//TODO: Particle gen

	c.SetCD(core.ActionSkill, eCD)

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

	// c.S.Status["heizouburst"] = c.Core.F + count
	c.Core.Status.AddStatus("heizouburst", duration)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Fudou Style Vacuum Slugger",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Anemo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	//TODO: does heizou burst snapshot?
	snap := c.Snapshot(&ai)

	// cb := func(a core.AttackCB) {

	// }
	c.AddTask(func() {
		for i, t := range c.Core.Targets {
			// skip non-enemy targets
			if t.Type() != core.TargettableEnemy {
				continue
			}
			if c.Base.Cons >= 4 {
				c.c4(i)
			}
			if i > 4 {
				break
			}

			c.irisDmg("Windmuster Iris", t)
		}
	}, "AuraCheck", f)

	c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(5, false, core.TargettableEnemy), f)

	//TODO: Check CD with or without delay, check energy consume frame
	c.SetCD(core.ActionBurst, 720)
	c.ConsumeEnergy(21)
	return f, a
}

//The first Windmuster Iris explosion in each Windmuster Kick will regenerate 9 Elemental Energy for Shikanoin Heizou.
//Every subsequent explosion in that Windmuster Kick will each regenerate an additional 1.5 Energy for Heizou.
//One Windmuster Kick can regenerate a total of 13.5 Energy for Heizou in this manner.
func (c *char) c4(i int) {
	energy := 0.0
	switch i {
	case 1:
		energy += 9.5
	case 2, 3:
		energy += 1.5
	case 4:
		energy += 1.0
	}
	c.AddEnergy("heizou c4", energy)
}

//Each Declension stack will increase the CRIT Rate of the Heartstopper Strike unleashed by 4%.
//When Heizou possesses Conviction, this Heartstopper Strike's CRIT DMG is increased by 32%.

func (c *char) c6() {
	//TODO: predamagemod is the way to go wiyth this? not sure but I cried about this enough already :C
	val := make([]float64, core.EndStatType)
	c.AddPreDamageMod(core.PreDamageMod{
		Key:    "heizou-c6",
		Expiry: c.Core.F + 600,
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			if atk.Info.AttackTag != core.AttackTagElementalArt && atk.Info.AttackTag != core.AttackTagElementalArtHold {
				return nil, false
			}
			val[core.CR] = float64(c.decStack) * 0.04
			if c.decStack >= 4 {
				val[core.CD] = 0.32
			}
			return val, true
		},
	})

}
