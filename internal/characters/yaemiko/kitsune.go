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
		Pattern: core.NewDefCircHit(2, false, core.TargettableEnemy),
	}
	c.AddTask(c.kitsuneTick(k), "start kitsune-tick", 30)
	if len(c.kitsunes) < 3 {
		//FIFO
		c.kitsunes = append(c.kitsunes, k)
	} else {
		//FIFO pop, popped kitsunes handled in kitsuneTick fn
		c.kitsunes = append(c.kitsunes[1:], k)
	}
}

func (c *char) kitsuneBurst(ai core.AttackInfo, src int) {
	n := 0
	c.Core.Log.Debugw("sky kitsune thunderbolt", "frame", c.Core.F, "event", core.LogCharacterEvent, "src", src)
	snap := c.Snapshot(&ai)
	for i, v := range c.kitsunes {
		if v.src == src {
			v.ae.Snapshot = snap
			c.Core.Combat.QueueAttackEvent(&v.ae, 94+54+i*24) // starts 54 after burst hit and 24 frames consecutively after
		} else {
			c.kitsunes[n] = v
			n++
		}
	}

	c.kitsunes = c.kitsunes[:n]
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
