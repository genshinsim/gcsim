package character

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

type EnergyEventPayload struct {
	Receiver     *CharWrapper
	PreEnergy    float64
	Change       float64
	EnergySource string
}

func (e *EnergyEventPayload) IsEventPayload() {}

type ParticlePayload struct {
	Particle Particle
}

func (p *ParticlePayload) IsEventPayload() {}

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
	c.events.Emit(event.OnEnergyChange, &EnergyEventPayload{
		Receiver:     c,
		PreEnergy:    preEnergy,
		Change:       e,
		EnergySource: src,
	})
	c.log.NewEvent("adding energy", glog.LogEnergyEvent, c.Index).
		Write("rec'd", e).
		Write("pre_recovery", preEnergy).
		Write("post_recovery", c.Energy).
		Write("source", src).
		Write("max_energy", c.EnergyMax)

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

	c.events.Emit(event.OnEnergyChange, &EnergyEventPayload{
		Receiver:     c,
		PreEnergy:    pre,
		Change:       amt,
		EnergySource: p.Source,
	})
	c.log.NewEvent(
		"particle",
		glog.LogEnergyEvent,
		c.Index,
	).
		Write("source", p.Source).
		Write("count", p.Num).
		Write("ele", p.Ele).
		Write("ER", er).
		Write("is_active", isActive).
		Write("party_count", partyCount).
		Write("pre_recovery", pre).
		Write("amt", amt).
		Write("post_recovery", c.Energy).
		Write("max_energy", c.EnergyMax)

}
