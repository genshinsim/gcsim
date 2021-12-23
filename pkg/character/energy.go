package character

import "github.com/genshinsim/gcsim/pkg/core"

func (c *Tmpl) QueueParticle(src string, num int, ele core.EleType, delay int) {
	p := core.Particle{
		Source: src,
		Num:    num,
		Ele:    ele,
	}
	c.AddTask(func() {
		c.Core.Energy.DistributeParticle(p)
	}, "particle", delay)

}

func (c *Tmpl) ConsumeEnergy(delay int) {
	if delay == 0 {
		c.Energy = 0
		return
	}
	c.AddTask(func() {
		c.Energy = 0
	}, "consume-energy", delay)
}

func (c *Tmpl) CurrentEnergy() float64 {
	return c.Energy
}

func (c *Tmpl) MaxEnergy() float64 {
	return c.EnergyMax
}

func (c *Tmpl) AddEnergy(e float64) {
	c.Energy += e
	if c.Energy > c.EnergyMax {
		c.Energy = c.EnergyMax
	}
	if c.Energy < 0 {
		c.Energy = 0
	}
	c.Core.Log.Debugw("adding energy", "frame", c.Core.F, "event", core.LogEnergyEvent, "rec'd", e, "next energy", c.Energy, "char", c.Index)
}

func (c *Tmpl) ReceiveParticle(p core.Particle, isActive bool, partyCount int) {
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
	case p.Ele == core.NoElement:
		amt = 2
	default:
		amt = 1
	}
	amt = amt * r //apply off field reduction
	//apply energy regen stat

	energyCalcModeBit := int8(0)
	if c.Core.Flags.EnergyCalcMode {
		energyCalcModeBit = 1
	}

	er = c.Stat(core.ER)

	amt = amt * (1 + er*(1-float64(energyCalcModeBit))) * float64(p.Num)

	pre := c.Energy

	c.Energy += amt
	if c.Energy > c.EnergyMax {
		c.Energy = c.EnergyMax
	}

	c.Core.Log.Debugw(
		"particle",
		"frame", c.Core.F,
		"event", core.LogEnergyEvent,
		"char", c.Index,
		"source", p.Source,
		"count", p.Num,
		"ele", p.Ele,
		"ER", er,
		"is_active", isActive,
		"party_count", partyCount,
		"pre_recovery", pre,
		"amt", amt,
		"post_recovery", c.Energy,
	)
}
