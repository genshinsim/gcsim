package heizou

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) skillHoldDuration(stacks int) int {
	//animation duration only
	//diff is the number of stacks we must charge up to reach the desired state
	diff := stacks - c.decStack
	if diff < 0 {
		diff = c.decStack
	}
	if diff > 4 {
		diff = 4
	}
	//it's .75s per stack
	return 45 * diff
}

func (c *char) addDecStack() {
	if c.decStack < 4 {
		c.decStack++
		c.Core.Log.NewEvent(
			"declension stack gained",
			core.LogCharacterEvent,
			c.Index,
			"stacks", c.decStack,
		)
	}
}

func (c *char) resetDecStack() {
	c.decStack = 0
}

const skillChargeStart = 12

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)
	stacks := p["hold"]
	dur := c.skillHoldDuration(stacks) //this should max out to 3s

	//queue task to increase stacks every 0.75s up to dur
	for i := 45; i <= dur; i++ {
		c.Core.Tasks.Add(func() {
			c.addDecStack()
		}, skillChargeStart+i)
	}

	//queue the attack as a task that goes through at the end of the animation; check for stacks then
	//animation should be skillChargeStart + dur + attack animation length
	c.Core.Tasks.Add(func() {
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
		if c.decStack == 4 {
			ai.Mult += convicBonus[c.TalentLvlSkill()]
		}
		//generate snap
		snap := c.Snapshot(&ai)
		//check for c6, increase crit
		if c.Base.Cons >= 6 {
			c.c6(&snap)
		}
		//a4
		done := false
		a4cb := func(a core.AttackCB) {
			if done {
				return
			}
			done = true
			c.a4()
		}
		//particle delayed 100 after dmg
		count := 2
		switch c.decStack {
		case 2, 3:
			if c.Core.Rand.Float64() < .5 {
				count++
			}
		case 4:
			count++
		}
		c.QueueParticle("heizou", count, core.Anemo, 100)
		//ok to reset stacks now
		c.Core.Log.NewEvent(
			"stack removed",
			core.LogCharacterEvent,
			c.Index,
		)
		c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(3, false, core.TargettableEnemy), 0, a4cb)
	}, skillChargeStart+dur+f)
	//TODO: Verify attack frame

	c.SetCD(core.ActionSkill, eCD)

	return skillChargeStart + dur + f, a
}
