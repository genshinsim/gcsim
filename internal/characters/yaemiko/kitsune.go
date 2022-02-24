package yaemiko

import (
	"log"

	"github.com/genshinsim/gcsim/pkg/core"
)

type kitsune struct {
	src     int
	deleted bool
}

func (c *char) makeKitsune() {
	k := &kitsune{}
	k.src = c.Core.F
	k.deleted = false
	//start ticking
	c.AddTask(c.kitsuneTick(k), "kitsune-tick", 45)
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
		c.Core.Status.AddStatus(yaeTotemStatus, 866)
	}
	//pop oldest first
	if len(c.kitsunes) == 3 {
		c.popOldestKitsune()
	}
	c.kitsunes = append(c.kitsunes, k)
	c.AddTag(yaeTotemCount, c.sakuraLevelCheck())

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
		dur := c.Core.F - c.kitsunes[0].src + 866
		if dur < 0 {
			log.Panicf("oldest totem should have expired already? dur: %v totem: %v", dur, *c.kitsunes[0])
		}
		c.Core.Status.AddStatus(yaeTotemStatus, dur)
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
		c.Core.Log.NewEvent("sky kitsune thunderbolt", core.LogCharacterEvent, c.Index, "src", c.kitsunes[i].src, "delay", 94+54+i*24)
	}
	// c.AddTask(func() {
	// 	//pop all?
	for range c.kitsunes {
		c.popOldestKitsune()
	}
	// }, "delay despawn for kitsunes", 0)

}

func (c *char) kitsuneTick(totem *kitsune) func() {

	return func() {
		//if deleted do nothing
		if totem.deleted {
			return
		}
		// c6
		// Sesshou Sakura start at Level 2 when created. Max level increased to 4, and their attacks will ignore 45% of the opponents' DEF.
		ai := core.AttackInfo{
			Abil:       "Sesshou Sakura Tick",
			ActorIndex: c.Index,
			AttackTag:  core.AttackTagElementalArt,
			ICDTag:     core.ICDTagElementalArt,
			ICDGroup:   core.ICDGroupDefault,
			StrikeType: core.StrikeTypeDefault,
			Element:    core.Electro,
			Durability: 25,
		}
		if c.Base.Cons >= 6 {
			ai.IgnoreDefPercent = 0.60
		}
		// no snapshot
		c.AddTask(func() {
			c.Core.Log.NewEvent("sky kitsune tick at level", core.LogCharacterEvent, c.Index, "sakura level", c.sakuraLevelCheck())
			ai.Mult = skill[c.sakuraLevelCheck()+c.turretBonus-1][c.TalentLvlSkill()]
			c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), 1, 1)
		}, "kitsune-tick", 50)
		if c.Core.F+51 >= c.totemLastParticleF+176 { // 176 frame ICD until we are sure about ICD
			c.Core.Log.NewEvent("sky kitsune particle", core.LogCharacterEvent, c.Index, "lastParticleF", c.totemLastParticleF)
			c.QueueParticle("yaemiko", 1, core.Electro, 51+30)
			c.totemLastParticleF = c.Core.F + 51
		}
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
