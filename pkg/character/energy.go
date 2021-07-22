package character

import "github.com/genshinsim/gsim/pkg/def"

func (c *Tmpl) QueueParticle(src string, num int, ele def.EleType, delay int) {
	p := def.Particle{
		Source: src,
		Num:    num,
		Ele:    ele,
	}
	c.AddTask(func() {
		c.Sim.DistributeParticle(p)
	}, "particle", delay)

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
	c.Log.Debugw("adding energy", "frame", c.Sim.Frame(), "event", def.LogEnergyEvent, "rec'd", e, "next energy", c.Energy, "char", c.Index)
}

func (c *Tmpl) ReceiveParticle(p def.Particle, isActive bool, partyCount int) {
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
	case p.Ele == def.NoElement:
		amt = 2
	default:
		amt = 1
	}
	amt = amt * r //apply off field reduction
	//apply energy regen stat
	er = c.Stat(def.ER)
	amt = amt * (1 + er) * float64(p.Num)

	pre := c.Energy

	c.Energy += amt
	if c.Energy > c.EnergyMax {
		c.Energy = c.EnergyMax
	}

	c.Log.Debugw(
		"particle",
		"frame", c.Sim.Frame(),
		"event", def.LogEnergyEvent,
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
