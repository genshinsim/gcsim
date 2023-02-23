package xingqiu

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

// start a new orbital or extended if already active; duration is length
// and delay is first tick starting
func (c *char) applyOrbital(duration int, delay int) {
	src := c.Core.F
	c.Core.Log.NewEvent(
		"Applying orbital", glog.LogCharacterEvent, c.Index,
	).Write(
		"current status", c.StatusExpiry(orbitalKey),
	)
	//check if orbitals already active, if active extend duration
	//other wise start first tick func
	if !c.orbitalActive {
		//use hitlag affected queue
		c.QueueCharTask(c.orbitalTickTask(src), delay)
		c.orbitalActive = true
		c.Core.Log.NewEvent(
			"orbital applied", glog.LogCharacterEvent, c.Index,
		).Write(
			"expected end", src+900,
		).Write(
			"next expected tick", src+40,
		)
	}
	c.AddStatus(orbitalKey, duration, true)
	c.Core.Log.NewEvent(
		"orbital duration extended", glog.LogCharacterEvent, c.Index,
	).Write(
		"new expiry", c.StatusExpiry(orbitalKey),
	)
}

func (c *char) orbitalTickTask(src int) func() {
	return func() {
		c.Core.Log.NewEvent(
			"orbital checking tick", glog.LogCharacterEvent, c.Index,
		).Write(
			"expiry", c.StatusExpiry(orbitalKey),
		).Write(
			"src", src,
		)
		if !c.StatusIsActive(orbitalKey) {
			c.orbitalActive = false
			return
		}

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Xingqiu Orbital",
			AttackTag:  attacks.AttackTagNone,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Hydro,
			Durability: 25,
		}
		c.Core.Log.NewEvent(
			"orbital ticked", glog.LogCharacterEvent, c.Index,
		).Write(
			"next expected tick", c.Core.F+135,
		).Write(
			"expiry", c.StatusExpiry(orbitalKey),
		).Write(
			"src", src,
		)

		//queue up next instance
		c.QueueCharTask(c.orbitalTickTask(src), 135)

		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 1.2), -1, 1)
	}
}
