package amber

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

type bunny struct {
	ae  combat.AttackEvent
	src int
}

// TODO: forbidden bunny cryo swirl tech
func (c *char) makeBunny() {
	b := bunny{}
	b.src = c.Core.F
	ai := combat.AttackInfo{
		Abil:       "Baron Bunny",
		ActorIndex: c.Index,
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Pyro,
		Durability: 50,
		Mult:       bunnyExplode[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)
	b.ae = combat.AttackEvent{
		Info:        ai,
		Pattern:     combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3),
		SourceFrame: c.Core.F,
		Snapshot:    snap,
	}

	c.bunnies = append(c.bunnies, b)

	//ondeath explodes
	//duration is 8.2 sec
	c.Core.Tasks.Add(func() {
		c.explode(b.src)
	}, 492)
}

func (c *char) explode(src int) {
	n := 0
	c.Core.Log.NewEvent("amber exploding bunny", glog.LogCharacterEvent, c.Index).
		Write("src", src)
	for _, v := range c.bunnies {
		if v.src == src {
			c.Core.QueueAttackEvent(&v.ae, 1)
			//4 orbs
			c.Core.QueueParticle("amber", 4, attributes.Pyro, c.ParticleDelay)
		} else {
			c.bunnies[n] = v
			n++
		}
	}

	c.bunnies = c.bunnies[:n]
}

func (c *char) manualExplode() {
	//do nothing if there are no bunnies
	if len(c.bunnies) == 0 {
		return
	}
	//only explode the first bunny
	if len(c.bunnies) > 0 {
		c.bunnies[0].ae.Info.Mult += 2
		c.Core.QueueAttackEvent(&c.bunnies[0].ae, 1)
		c.Core.QueueParticle("amber", 4, attributes.Pyro, c.ParticleDelay)
	}
	c.bunnies = c.bunnies[1:]
}

// explode all bunnies on overload
func (c *char) overloadExplode() {
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {

		atk := args[1].(*combat.AttackEvent)
		if len(c.bunnies) == 0 {
			return false
		}
		//TODO: only amber trigger?
		if atk.Info.ActorIndex != c.Index {
			return false
		}

		if atk.Info.AttackTag != combat.AttackTagOverloadDamage {
			return false
		}

		for _, v := range c.bunnies {
			//every bunny gets bonus multiplikers
			v.ae.Info.Mult += 2
			c.Core.QueueAttackEvent(&v.ae, 1)
			c.Core.QueueParticle("amber", 4, attributes.Pyro, c.ParticleDelay)
		}
		c.bunnies = make([]bunny, 0, 2)

		return false
	}, "amber-bunny-explode-overload")
}
