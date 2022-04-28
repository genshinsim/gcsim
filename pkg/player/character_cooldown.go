package player

import "github.com/genshinsim/gcsim/pkg/core/glog"

type CooldownModFunc func(a Action) float64

type cooldownMod struct {
	Key    string
	Amount func(a Action) float64
	Expiry int
}

func (c *MasterChar) CDReduction(a Action, dur int) int {
	var cd float64 = 1
	n := 0
	for _, v := range c.CooldownMods {
		//if not expired
		if v.Expiry == -1 || v.Expiry > c.Player.Core.F {
			amt := v.Amount(a)
			c.Player.Core.Log.NewEvent(
				"applying cooldown modifier",
				glog.LogActionEvent,
				c.Index,
				"key", v.Key,
				"modifier", amt,
				"expiry", v.Expiry,
			)
			cd += amt
			c.CooldownMods[n] = v
			n++
		}
	}
	c.CooldownMods = c.CooldownMods[:n]

	return int(float64(dur) * cd)
}

func (c *MasterChar) AddCDAdjust(key string, amt func(a Action) float64, duration int) {
	ind := -1
	for i, v := range c.CooldownMods {
		//if expired already, set to nil and ignore
		if v.Key == key {
			ind = i
		}
	}
	if ind > -1 {
		c.CooldownMods[ind] = cooldownMod{
			Key:    key,
			Amount: amt,
			Expiry: c.Player.Core.F + duration,
		}
	} else {
		c.CooldownMods = append(c.CooldownMods, cooldownMod{
			Key:    key,
			Amount: amt,
			Expiry: c.Player.Core.F + duration,
		})
	}
}
