package amber

import "github.com/genshinsim/gcsim/pkg/core"

type bunny struct {
	ae  core.AttackEvent
	src int
}

//TODO: forbidden bunny cryo swirl tech
func (c *char) makeBunny() {
	b := bunny{}
	b.src = c.Core.F
	ai := core.AttackInfo{
		Abil:       "Baron Bunny",
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Pyro,
		Durability: 50,
		Mult:       bunnyExplode[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)
	b.ae = core.AttackEvent{
		Info:        ai,
		Pattern:     core.NewDefCircHit(2, false, core.TargettableEnemy),
		SourceFrame: c.Core.F,
		Snapshot:    snap,
	}

	c.bunnies = append(c.bunnies, b)

	//ondeath explodes
	//duration is 8.2 sec
	c.AddTask(func() {
		c.explode(b.src)
	}, "bunny", 492)
}

func (c *char) explode(src int) {
	n := 0
	c.Core.Log.Debugw("amber exploding bunny", "frame", c.Core.F, "event", core.LogCharacterEvent, "src", src)
	for _, v := range c.bunnies {
		if v.src == src {
			c.Core.Combat.QueueAttackEvent(&v.ae, 1)
			//4 orbs
			c.QueueParticle("amber", 4, core.Pyro, 100)
		} else {
			c.bunnies[n] = v
			n++
		}
	}

	c.bunnies = c.bunnies[:n]
}

func (c *char) manualExplode() {
	//only explode the first bunny
	if len(c.bunnies) > 0 {
		c.bunnies[0].ae.Info.Mult += 2
		c.Core.Combat.QueueAttackEvent(&c.bunnies[0].ae, 1)
		c.QueueParticle("amber", 4, core.Pyro, 100)
	}
	c.bunnies = c.bunnies[1:]
}

func (c *char) overloadExplode() {
	//explode all bunnies on overload

	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {

		atk := args[1].(*core.AttackEvent)
		if len(c.bunnies) == 0 {
			return false
		}
		//TODO: only amber trigger?
		if atk.Info.ActorIndex != c.Index {
			return false
		}

		if atk.Info.AttackTag != core.AttackTagOverloadDamage {
			return false
		}

		for _, v := range c.bunnies {
			//every bunny gets bonus multiplikers
			v.ae.Info.Mult += 2
			c.Core.Combat.QueueAttackEvent(&v.ae, 1)
			c.QueueParticle("amber", 4, core.Pyro, 100)
		}
		c.bunnies = make([]bunny, 0, 2)

		return false
	}, "amber-bunny-explode-overload")

}
