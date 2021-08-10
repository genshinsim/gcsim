package character

import (
	"github.com/genshinsim/gsim/pkg/core"
)

func (c *Tmpl) Attack(p map[string]int) int {
	return 0
}

func (c *Tmpl) Aimed(p map[string]int) int {
	return 0
}

func (c *Tmpl) ChargeAttack(p map[string]int) int {
	return 0
}

func (c *Tmpl) HighPlungeAttack(p map[string]int) int {
	return 0
}

func (c *Tmpl) LowPlungeAttack(p map[string]int) int {
	return 0
}

func (c *Tmpl) Skill(p map[string]int) int {
	return 0
}

func (c *Tmpl) Burst(p map[string]int) int {
	return 0
}

func (c *Tmpl) Dash(p map[string]int) int {
	return 24
}

func (c *Tmpl) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	default:
		c.Log.Warnw("ActionStam not implemented", "character", c.Base.Name)
		return 0
	}
}

func (c *Tmpl) ActionFrames(a core.ActionType, p map[string]int) int {
	c.Log.Warnw("ActionFrames not implemented", "character", c.Base.Name)
	return 0
}

func (c *Tmpl) ActionReady(a core.ActionType, p map[string]int) bool {
	switch a {
	case core.ActionBurst:
		if c.Energy != c.EnergyMax {
			c.Log.Warnw("burst not enough energy")
			return false
		}
		return c.ActionCD[a] <= c.Sim.Frame()
	case core.ActionSkill:
		return c.ActionCD[a] <= c.Sim.Frame()
	}
	return true
}

func (c *Tmpl) AddCDAdjustFunc(rd core.CDAdjust) {
	ind := len(c.CDReductionFuncs)
	for i, v := range c.CDReductionFuncs {
		//if expired already, set to nil and ignore
		if v.Key == rd.Key {
			ind = i
		}
	}
	if ind == len(c.CDReductionFuncs) {
		c.CDReductionFuncs = append(c.CDReductionFuncs, rd)
	} else {
		c.CDReductionFuncs[ind] = rd
	}
}

func (c *Tmpl) SetCD(a core.ActionType, dur int) {
	//here we reduce dur by cd reduction
	var cd float64 = 1
	n := 0
	for _, v := range c.CDReductionFuncs {
		//if not expired
		if v.Expiry == -1 || v.Expiry > c.Sim.Frame() {
			amt := v.Amount(a)
			c.Log.Debugw("applying cooldown modifier", "frame", c.Sim.Frame(), "event", core.LogCharacterEvent, "char", c.Index, "key", v.Key, "modifier", amt, "expiry", v.Expiry)
			cd += amt
			c.CDReductionFuncs[n] = v
			n++
		}
	}
	c.CDReductionFuncs = c.CDReductionFuncs[:n]

	c.ActionCD[a] = c.Sim.Frame() + int(float64(dur)*cd) //truncate to int
	c.Log.Debugw("cooldown triggered", "frame", c.Sim.Frame(), "event", core.LogCharacterEvent, "char", c.Index, "type", a.String(), "expiry", c.Sim.Frame()+dur)
}

func (c *Tmpl) Cooldown(a core.ActionType) int {
	cd := c.ActionCD[a] - c.Sim.Frame()
	if cd < 0 {
		cd = 0
	}
	return cd
}

func (c *Tmpl) ResetActionCooldown(a core.ActionType) {
	c.ActionCD[a] = 0
}

func (c *Tmpl) ReduceActionCooldown(a core.ActionType, v int) {
	c.ActionCD[a] -= v
}

func (c *Tmpl) ResetNormalCounter() {
	c.NormalCounter = 0
}
