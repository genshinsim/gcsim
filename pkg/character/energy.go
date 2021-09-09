package character

import "github.com/genshinsim/gsim/pkg/core"

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

func (c *Tmpl) ConsumeEnergy(refund float64, delay int) {
	c.ER_CalcValues()
	if delay > 0 { //not sure how this task delay interacs whit the energy clear and ER calcs on a 0 delay, so forcing it to no interact
		c.AddTask(func() {
			c.Energy = refund
			c.ER_ClearEnergy()
		}, "consume-energy", delay)
	} else {
		c.Energy = refund
		c.ER_ClearEnergy()
	}
}

func (c *Tmpl) CurrentEnergy() float64 {
	return c.Energy
}

func (c *Tmpl) MaxEnergy() float64 {
	return c.EnergyMax
}

func (c *Tmpl) AddEnergy(e float64) {
	c.Energy += e
	c.ER_SaveFlatEnergy(e) //function call ER calcs, e is the energy added to a character
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
	er = c.Stat(core.ER)

	c.ER_SaveEnergy(amt, p) //function call ER calcs, amt must be pre char ER multiplication and after off field reduction, p is particle recived

	amt = amt * (1 + er) * float64(p.Num)

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
		// addign ER text for debug and test before creating is own log category
		"ER energy", c.ER_EnergyRecived,
		"ER flat energy", c.ER_FlatEnergyRecived,
		"ER need", c.ER_ERneeded-1, //-1 to be consistent with the rest of the log
	)
}

func (c *Tmpl) ER_SaveEnergy(amt float64, p core.Particle) { //saves the particles recived in a char pre current ER multiplication, only for a one Burst
	if !(c.Energy == c.EnergyMax && c.ER_EnergyRecived == 0 && c.ER_FlatEnergyRecived == 0) { //workaround, dont save when the burst is already full at beginning of simulation
		c.ER_EnergyRecived += amt * float64(p.Num)
	}
}

func (c *Tmpl) ER_SaveFlatEnergy(amt float64) { //saves the falt particles recived in a char (3 energy from c6 xq), only for a one Burst
	c.ER_FlatEnergyRecived += amt
}

func (c *Tmpl) ER_ClearEnergy() { //Restart the ER especific values wen a burst is used
	c.ER_EnergyRecived = 0
	c.ER_FlatEnergyRecived = 0
}

func (c *Tmpl) ER_CalcValues() { //Temporary max ER value calcs and storage in a variable because dont know how to log and display as a individual category
	var ERneed = (c.Energy - c.ER_FlatEnergyRecived) / c.ER_EnergyRecived
	if c.ER_EnergyRecived > 0 { //workaround for using filling burst at start of simulation
		if c.ER_ERneeded == 0 { //chechking for the highest value of ER and storing it while the log and display of functions are not finished
			c.ER_ERneeded = ERneed
		} else {
			if c.ER_ERneeded < ERneed {
				c.ER_ERneeded = ERneed
			}
		}
	}
}
