package xingqiu

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

//start a new orbital or extended if already active; duration is length
//and delay is first tick starting
func (c *char) applyOrbital(duration int, delay int) {
	src := c.Core.F
	c.Core.Log.NewEvent(
		"Applying orbital", glog.LogCharacterEvent, c.Index,
		"current status", c.StatusExpiry(orbitalKey),
	)
	//check if orbitals already active, if active extend duration
	//other wise start first tick func
	if !c.orbitalActive {
		c.Core.Tasks.Add(c.orbitalTickTask(src), delay)
		c.orbitalActive = true
		c.Core.Log.NewEvent(
			"orbital applied", glog.LogCharacterEvent, c.Index,
			"expected end", src+900,
			"next expected tick", src+40,
		)
	}
	c.AddStatus(orbitalKey, duration, true)
	c.Core.Log.NewEvent(
		"orbital duration extended", glog.LogCharacterEvent, c.Index,
		"new expiry", c.StatusExpiry(orbitalKey),
	)
}

func (c *char) orbitalTickTask(src int) func() {
	return func() {
		c.Core.Log.NewEvent(
			"orbital checking tick", glog.LogCharacterEvent, c.Index,
			"expiry", c.StatusExpiry(orbitalKey),
			"src", src,
		)
		if !c.StatusIsActive(orbitalKey) {
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
		c.Core.Log.NewEvent(
			"orbital ticked", glog.LogCharacterEvent, c.Index,
			"next expected tick", c.Core.F+135,
			"expiry", c.StatusExpiry(orbitalKey),
			"src", src,
		)

		//queue up next instance
		c.Core.Tasks.Add(c.orbitalTickTask(src), 135)

		c.Core.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), -1, 1)
	}
}
