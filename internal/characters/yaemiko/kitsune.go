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
	c.AddTask(c.kitsuneTick(k), "start kitsune-tick", 60)
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
			c.Core.Combat.QueueAttackEvent(&v.ae, 30+i*30)
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
			Mult:       skill[c.sakuraLevelCheck()+c.c6int][c.TalentLvlSkill()],
		}
		if c.Base.Cons >= 6 {
			ai.IgnoreDefPercent = 0.45
		}
		c.Core.Log.Debugw("sky kitsune tick", "frame", c.Core.F, "event", core.LogCharacterEvent)

		c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), -1, 5)

		// tick per 2.5 seconds
		c.AddTask(c.kitsuneTick(totem), "kitsune-tick", 150)
	}
}

func (c *char) sakuraLevelCheck() int {
	if len(c.kitsunes) > 3 {
		panic("wtf more than 3 kitsunes")
	} else {
		return len(c.kitsunes)
	}
}
