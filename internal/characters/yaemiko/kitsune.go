package yaemiko

import "github.com/genshinsim/gcsim/pkg/core"

type kitsune struct {
	ae  core.AttackEvent
	src int
}

func (c *char) makeKitsune() {
	k := kitsune{}
	k.src = c.Core.F
	ai := core.AttackInfo{
		Abil:       "Sky Kitsune Thunderbolt",
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Electro,
		Durability: 50,
		Mult:       burst[1][c.TalentLvlSkill()],
	}
	k.ae = core.AttackEvent{
		Info:    ai,
		Pattern: core.NewDefCircHit(5, false, core.TargettableEnemy),
	}
	c.AddTask(c.kitsuneTick(k), "start kitsune-tick", 30)
	if len(c.kitsunes) < 3 {
		//FIFO
		c.kitsunes = append(c.kitsunes, k)
		c.Tags["totems"]++
	} else {
		//FIFO pop, popped kitsunes handled in kitsuneTick fn
		c.kitsunes = append(c.kitsunes[1:], k)
	}
	if len(c.kitsunes) == 0 {
		c.Core.Status.AddStatus("oldestTotemExpiry", 14*60)
	}
}

func (c *char) kitsuneBurst(ai core.AttackInfo, sakuraLevel int) {
	snap := c.Snapshot(&ai)
	for i := 0; i < sakuraLevel; i++ {
		c.kitsunes[i].ae.Snapshot = snap
		c.Core.Combat.QueueAttackEvent(&c.kitsunes[i].ae, 94+54+i*24) // starts 54 after burst hit and 24 frames consecutively after
		c.Core.Log.Debugw("sky kitsune thunderbolt", "frame", c.Core.F, "event", core.LogCharacterEvent, "src", c.kitsunes[i].src, "delay", 94+54+i*24)
	}
	c.AddTask(func() {
		c.kitsunes = c.kitsunes[:0]
		c.Tags["totems"] = 0
		c.Core.Status.DeleteStatus("oldestTotemExpiry")
	}, "delay despawn for kitsunes", 78)

}

func (c *char) kitsuneTick(totem kitsune) func() {

	return func() {
		//make sure it's not overwritten
		flag := false
		for _, v := range c.kitsunes {
			if v.src == totem.src {
				if flag {
					panic("two kitsune's found created at the same time")
				}
				flag = true
			}
		}
		if !flag {
			// don't perform kitsune tick if kitsune does not exist
			return
		}
		//do nothing if totem expired
		if c.Core.F > totem.src+14*60 {
			c.Tags["totems"]--
			// c.kitsunes = c.kitsunes[1:]
			// if len(c.kitsunes) > 0 {
			// 	if c.kitsunes[0].src+14*60-c.Core.F > 0 {
			// 		c.Core.Status.AddStatus("oldestTotemExpiry", c.cdQueue[core.ActionSkill][0])
			// 	}
			// }
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
			Mult:       skill[c.sakuraLevelCheck()+c.c6int-1][c.TalentLvlSkill()],
		}
		if c.Base.Cons >= 6 {
			ai.IgnoreDefPercent = 0.45
		}
		c.Core.Log.Debugw("sky kitsune tick", "frame", c.Core.F, "event", core.LogCharacterEvent)
		// no snapshot
		c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), 0, 49)
		if c.Core.F+49 >= c.totemLastParticleF+60*2.5 {
			c.Core.Log.Debugw("sky kitsune particle", "frame", c.Core.F, "event", core.LogCharacterEvent, "lastParticleF", c.totemLastParticleF)
			c.QueueParticle("kitsune-tick particle", 1, core.Electro, 49+30)
			c.totemLastParticleF = c.Core.F + 49
		}
		// tick per 2.5 seconds
		c.AddTask(c.kitsuneTick(totem), "kitsune-tick", 177)
	}
}

func (c *char) sakuraLevelCheck() int {
	sakuraLevels := 0
	for _, v := range c.kitsunes {
		if c.Core.F < v.src+14*60 {
			sakuraLevels++
		}
	}
	if sakuraLevels > 3 {
		panic("wtf more than 3 kitsunes")
	} else {
		return sakuraLevels
	}
}
