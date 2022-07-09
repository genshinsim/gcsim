package xingqiu

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (c *char) applyOrbital(duration int, delay int) {
	src := c.Core.F
	c.Core.Log.NewEvent("Applying orbital", glog.LogCharacterEvent, c.Index).
		Write("current status", c.Core.Status.Duration("xqorb"))
	//check if orbitals already active, if active extend duration
	//other wise start first tick func
	if !c.orbitalActive {
		c.Core.Tasks.Add(c.orbitalTickTask(src), delay)
		c.orbitalActive = true
		c.Core.Log.NewEvent("orbital applied", glog.LogCharacterEvent, c.Index).
			Write("expected end", src+900).
			Write("next expected tick", src+40)
	}

	c.Core.Status.Add("xqorb", duration)
	c.Core.Log.NewEvent("orbital duration extended", glog.LogCharacterEvent, c.Index).
		Write("new expiry", c.Core.Status.Duration("xqorb"))
}

func (c *char) orbitalTickTask(src int) func() {
	return func() {
		c.Core.Log.NewEvent("orbital checking tick", glog.LogCharacterEvent, c.Index).
			Write("expiry", c.Core.Status.Duration("xqorb")).
			Write("src", src)
		if c.Core.Status.Duration("xqorb") == 0 {
			c.orbitalActive = false
			return
		}

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Xingqiu Orbital",
			AttackTag:  combat.AttackTagNone,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			Element:    attributes.Hydro,
			Durability: 25,
		}
		c.Core.Log.NewEvent("orbital ticked", glog.LogCharacterEvent, c.Index).
			Write("next expected tick", c.Core.F+135).
			Write("expiry", c.Core.Status.Duration("xqorb")).
			Write("src", src)

		//queue up next instance
		c.Core.Tasks.Add(c.orbitalTickTask(src), 135)

		c.Core.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), -1, 1)
	}
}
