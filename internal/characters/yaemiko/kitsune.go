package yaemiko

import (
	"log"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

type kitsune struct {
	src     int
	deleted bool
}

func (c *char) makeKitsune() {
	k := &kitsune{}
	k.src = c.Core.Frame
	k.deleted = false
	//start ticking
	c.AddTask(c.kitsuneTick(k), "kitsune-tick", 45+50)
	//add task to delete this one if times out (and not deleted by anything else)
	c.AddTask(func() {
		//i think we can just check for .deleted here
		if k.deleted {
			return
		}
		//ok now we can delete this
		c.popOldestKitsune()
	}, "kitsune-expiry", 866) // e ani + duration

	if len(c.kitsunes) == 0 {
		c.Core.AddStatus(yaeTotemStatus, 866)
	}
	//pop oldest first
	if len(c.kitsunes) == 3 {
		c.popOldestKitsune()
	}
	c.kitsunes = append(c.kitsunes, k)
	c.AddTag(yaeTotemCount, c.sakuraLevelCheck())

}

func (c *char) popAllKitsune() {
	for i := range c.kitsunes {
		c.kitsunes[i].deleted = true
	}
	c.kitsunes = c.kitsunes[:0]
	c.Core.Status.DeleteStatus(yaeTotemStatus)
	c.AddTag(yaeTotemCount, 0)
}

func (c *char) popOldestKitsune() {
	if len(c.kitsunes) == 0 {
		//nothing to pop??
		return
	}

	c.kitsunes[0].deleted = true
	c.kitsunes = c.kitsunes[1:]

	//here check for status
	if len(c.kitsunes) > 0 {
		dur := c.Core.Frame - c.kitsunes[0].src + 866
		if dur < 0 {
			log.Panicf("oldest totem should have expired already? dur: %v totem: %v", dur, *c.kitsunes[0])
		}
		c.Core.AddStatus(yaeTotemStatus, dur)
	} else {
		c.Core.Status.DeleteStatus(yaeTotemStatus)
	}

	c.AddTag(yaeTotemCount, len(c.kitsunes))
}

func (c *char) kitsuneBurst(ai core.AttackInfo, pattern core.AttackPattern) {
	for i := 0; i < c.sakuraLevelCheck(); i++ {
		c.Core.Combat.QueueAttack(ai, pattern, 94+54+i*24, 94+54+i*24) // starts 54 after burst hit and 24 frames consecutively after
		if c.Base.Cons >= 1 {
			c.AddTask(func() {
				c.AddEnergy("yae-c1", 8)
			}, "energy from sky kitsune", 94+54+i*24)
		}
		c.ResetActionCooldown(core.ActionSkill)
		c.coretype.Log.NewEvent("sky kitsune thunderbolt", coretype.LogCharacterEvent, c.Index, "src", c.kitsunes[i].src, "delay", 94+54+i*24)
	}
	// c.AddTask(func() {
	// 	//pop all?
	// for range c.kitsunes {
	// 	c.popOldestKitsune()
	// }
	// }, "delay despawn for kitsunes", 0)
	c.popAllKitsune()

}

func (c *char) kitsuneTick(totem *kitsune) func() {
	return func() {
		//if deleted do nothing
		if totem.deleted {
			return
		}
		// c6
		// Sesshou Sakura start at Level 2 when created. Max level increased to 4, and their attacks will ignore 45% of the opponents' DEF.

		lvl := c.sakuraLevelCheck() - 1
		if c.Base.Cons >= 2 {
			lvl += 1
		}

		ai := core.AttackInfo{
			Abil:       "Sesshou Sakura Tick",
			ActorIndex: c.Index,
			AttackTag:  core.AttackTagElementalArt,
			Mult:       skill[lvl][c.TalentLvlSkill()],
			ICDTag:     core.ICDTagElementalArt,
			ICDGroup:   core.ICDGroupDefault,
			StrikeType: core.StrikeTypeDefault,
			Element:    core.Electro,
			Durability: 25,
		}

		c.coretype.Log.NewEvent("sky kitsune tick at level", coretype.LogCharacterEvent, c.Index, "sakura level", lvl)

		if c.Base.Cons >= 6 {
			ai.IgnoreDefPercent = 0.60
		}

		done := false
		cb := func(ac core.AttackCB) {
			if c.Base.Cons >= 4 && !done {
				done = true
				c.c4()
			}

			//on hit check for particles
			c.coretype.Log.NewEvent("sky kitsune particle", coretype.LogCharacterEvent, c.Index, "lastParticleF", c.totemParticleICD)
			if c.Core.Frame < c.totemParticleICD {
				return
			}
			c.totemParticleICD = c.Core.Frame + 176
			c.QueueParticle("yaemiko", 1, core.Electro, 30)
		}

		c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(c.Core.RandomEnemyTarget(), core.TargettableEnemy), 1, 1, cb)
		// tick per 2.5 seconds
		c.AddTask(c.kitsuneTick(totem), "kitsune-tick", 176)
	}
}

func (c *char) sakuraLevelCheck() int {

	count := len(c.kitsunes)

	if count < 0 {
		//this is for the base case when there are no totems (other wise we'll end up with 1 if C6)
		return 0
	}
	if count > 3 {
		panic("wtf more than 3 totems")
	}
	return count
}
