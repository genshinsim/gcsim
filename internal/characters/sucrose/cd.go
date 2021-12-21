package sucrose

import "github.com/genshinsim/gcsim/pkg/core"

//use a tick base system for just sucrose skill

func (c *char) ActionReady(a core.ActionType, p map[string]int) bool {
	if a != core.ActionSkill {
		return c.Tmpl.ActionReady(a, p)
	}
	//useable as long as there's more than 1 charge
	return c.eCharges > 0
}

func (c *char) Tick() {
	c.Tmpl.Tick()

	//do nothing if already at 0
	if c.ActionCD[core.ActionSkill] == 0 {
		return
	}

	//otherwise tick down
	c.ReduceActionCooldown(core.ActionSkill, 1)
}

func (c *char) ResetActionCooldown(a core.ActionType) {
	if a != core.ActionSkill {
		c.Tmpl.ResetActionCooldown(a)
		return
	}

	//basically we want to reduce cooldown by the existing amt
	c.ReduceActionCooldown(core.ActionSkill, eCD)

}

func (c *char) ReduceActionCooldown(a core.ActionType, v int) {
	if a != core.ActionSkill {
		c.Tmpl.ReduceActionCooldown(a, v)
		return
	}

	c.ActionCD[a] -= v

	//if we hit 0, then add 1 to stack and reset
	//cd to max
	if c.ActionCD[core.ActionSkill] <= 0 {
		c.ActionCD[core.ActionSkill] = 0
		c.eCharges++
		if c.eCharges >= c.eChargeMax {
			//note that we shouldn't ever have a case where eCharge > eChargeMax here
			//set to equal just in case
			c.eCharges = c.eChargeMax
			return
		}
		//if we're not at max charge yet then queue up another cd
		c.ActionCD[core.ActionSkill] = c.calcSkillCD(eCD)
	}

}

func (c *char) Cooldown(a core.ActionType) int {
	if a != core.ActionSkill {
		return c.Tmpl.Cooldown(a)
	}

	return c.ActionCD[a]
}

func (c *char) calcSkillCD(dur int) int {
	//here we reduce dur by cd reduction
	var cd float64 = 1
	n := 0
	for _, v := range c.CDReductionFuncs {
		//if not expired
		if v.Expiry == -1 || v.Expiry > c.Core.F {
			amt := v.Amount(core.ActionSkill)
			c.Core.Log.Debugw("applying cooldown modifier", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "key", v.Key, "modifier", amt, "expiry", v.Expiry)
			cd += amt
			c.CDReductionFuncs[n] = v
			n++
		}
	}
	c.CDReductionFuncs = c.CDReductionFuncs[:n]

	return int(float64(dur) * cd)
}

func (c *char) SetCD(a core.ActionType, dur int) {
	if a != core.ActionSkill {
		c.Tmpl.SetCD(a, dur)
		return
	}
	//reduce charge by 1
	c.eCharges--
	if c.eCharges < 0 {
		panic("sucrose e charge < 0")
	}
	//start ticking if it wasn't ticking before; other wise do nothing since we're already ticking
	if c.ActionCD[a] == 0 {
		c.ActionCD[a] = c.calcSkillCD(dur)
	}

	c.Core.Log.Debugw("cooldown triggered", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "type", a.String(), "expiry", c.Core.F+dur)
}
