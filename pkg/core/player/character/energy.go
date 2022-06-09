package character

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (c *CharWrapper) ConsumeEnergy(delay int) {
	if delay == 0 {
		c.Energy = 0
		return
	}
	c.tasks.Add(func() {
		c.Energy = 0
	}, delay)
}

func (c *CharWrapper) AddEnergy(src string, e float64) {
	preEnergy := c.Energy
	c.Energy += e
	if c.Energy > c.EnergyMax {
		c.Energy = c.EnergyMax
	}
	if c.Energy < 0 {
		c.Energy = 0
	}

	c.events.Emit(event.OnEnergyChange, c, preEnergy, e, src)
	c.log.NewEvent("adding energy", glog.LogEnergyEvent, c.Index,
		"rec'd", e,
		"post_recovery", c.Energy,
		"source", src,
		"max_energy", c.EnergyMax,
	)
}

func (c *CharWrapper) ReceiveParticle(p Particle, isActive bool, partyCount int) {
	var amt, er, r float64
	r = 1.0
	if !isActive {
		r = 1.0 - 0.1*float64(partyCount)
	}
	//recharge amount - particles: same = 3, non-ele = 2, diff = 1
	//recharge amount - orbs: same = 9, non-ele = 6, diff = 3 (3x particles)
	switch {
	case p.Ele == c.Base.Element:
		amt = 3
	case p.Ele == attributes.NoElement:
		amt = 2
	default:
		amt = 1
	}
	amt = amt * r //apply off field reduction
	//apply energy regen stat

	er = c.Stat(attributes.ER)

	amt = amt * (1 + er) * float64(p.Num)

	pre := c.Energy

	c.Energy += amt
	if c.Energy > c.EnergyMax {
		c.Energy = c.EnergyMax
	}

	c.events.Emit(event.OnEnergyChange, c, pre, amt, p.Source)
	c.log.NewEvent(
		"particle",
		glog.LogEnergyEvent,
		c.Index,
		"source", p.Source,
		"count", p.Num,
		"ele", p.Ele,
		"ER", er,
		"is_active", isActive,
		"party_count", partyCount,
		"pre_recovery", pre,
		"amt", amt,
		"post_recovery", c.Energy,
		"max_energy", c.EnergyMax,
	)
}
