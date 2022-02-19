package character

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *Tmpl) Attack(p map[string]int) (int, int) {
	return 0, 0
}

func (c *Tmpl) Aimed(p map[string]int) (int, int) {
	return 0, 0
}

func (c *Tmpl) ChargeAttack(p map[string]int) (int, int) {
	return 0, 0
}

func (c *Tmpl) HighPlungeAttack(p map[string]int) (int, int) {
	return 0, 0
}

func (c *Tmpl) LowPlungeAttack(p map[string]int) (int, int) {
	return 0, 0
}

func (c *Tmpl) Skill(p map[string]int) (int, int) {
	return 0, 0
}

func (c *Tmpl) Burst(p map[string]int) (int, int) {
	return 0, 0
}

func (c *Tmpl) Dash(p map[string]int) (int, int) {
	return 24, 24
}

func (c *Tmpl) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	default:
		c.Core.Log.NewEvent("ActionStam not implemented", core.LogActionEvent, c.Index)
		return 0
	}
}

func (c *Tmpl) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	c.Core.Log.NewEvent("ActionFrames not implemented", core.LogActionEvent, c.Index)
	return 0, 0
}

func (c *Tmpl) ActionReady(a core.ActionType, p map[string]int) bool {
	switch a {
	case core.ActionBurst:
		if (c.Energy != c.EnergyMax) && !c.Core.Flags.EnergyCalcMode {
			// c.Core.Log.Warnw("burst not enough energy",  core.LogActionEvent, "char", c.Index)
			return false
		}
		return c.ActionCD[a] <= c.Core.F
	case core.ActionSkill:
		return c.ActionCD[a] <= c.Core.F
	}
	return true
}

func (c *Tmpl) ActionInterruptableDelay(next core.ActionType) int {
	return 0
}

func (c *Tmpl) AddCDAdjustFunc(rd core.CDAdjust) {
	ind := -1
	for i, v := range c.CDReductionFuncs {
		//if expired already, set to nil and ignore
		if v.Key == rd.Key {
			ind = i
		}
	}
	if ind > -1 {
		c.CDReductionFuncs[ind] = rd
	} else {
		c.CDReductionFuncs = append(c.CDReductionFuncs, rd)
	}
}

// Sets cooldown for a given action, applying any modifications
func (c *Tmpl) SetCD(a core.ActionType, dur int) {
	//here we reduce dur by cd reduction
	var cd float64 = 1
	n := 0
	for _, v := range c.CDReductionFuncs {
		//if not expired
		if v.Expiry == -1 || v.Expiry > c.Core.F {
			amt := v.Amount(a)
			c.Core.Log.NewEvent("applying cooldown modifier", core.LogActionEvent, c.Index, "key", v.Key, "modifier", amt, "expiry", v.Expiry)
			cd += amt
			c.CDReductionFuncs[n] = v
			n++
		}
	}
	c.CDReductionFuncs = c.CDReductionFuncs[:n]

	c.ActionCD[a] = c.Core.F + int(float64(dur)*cd) //truncate to int
	// Log to actions for the purpose of visibility since CDs are decently important
	c.Core.Log.NewEvent("cooldown triggered", core.LogActionEvent, c.Index, "type", a.String(), "expiry", c.Core.F+dur)
}

// Thin wrapper around SetCD to allow for setting CD after some delay frames
// Useful for bursts that consume energy and start CD at some point into their animation
func (c *Tmpl) SetCDWithDelay(a core.ActionType, dur int, delay int) {
	if delay == 0 {
		c.SetCD(a, dur)
		return
	}
	c.AddTask(func() { c.SetCD(a, dur) }, "set-cd", delay)
}

func (c *Tmpl) Cooldown(a core.ActionType) int {
	cd := c.ActionCD[a] - c.Core.F
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
