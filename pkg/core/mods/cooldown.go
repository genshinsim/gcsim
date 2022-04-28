package mods

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

type CooldownModFunc func(a action.Action) float64

type cooldownMod struct {
	Key    string
	Amount func(a action.Action) float64
	Expiry int
}

func (c *Handler) CDReduction(a action.Action, dur int, char int) int {
	var cd float64 = 1
	n := 0
	for _, v := range c.CooldownMods[char] {
		//if not expired
		if v.Expiry == -1 || v.Expiry > *c.f {
			amt := v.Amount(a)
			c.log.NewEvent(
				"applying cooldown modifier",
				glog.LogActionEvent,
				char,
				"key", v.Key,
				"modifier", amt,
				"expiry", v.Expiry,
			)
			cd += amt
			c.CooldownMods[char][n] = v
			n++
		}
	}
	c.CooldownMods[char] = c.CooldownMods[char][:n]

	return int(float64(dur) * cd)
}

func (c *Handler) AddCDAdjust(key string, dur int, amt func(a action.Action) float64, chars ...int) {
	for _, char := range chars {
		mod := cooldownMod{
			Key:    key,
			Amount: amt,
			Expiry: *c.f + dur,
		}
		ind := -1
		for i, v := range c.CooldownMods[char] {
			//if expired already, set to nil and ignore
			if v.Key == key {
				ind = i
			}
		}
		if ind > -1 {
			c.CooldownMods[char][ind] = mod
		} else {
			c.CooldownMods[char] = append(c.CooldownMods[char], mod)
		}
	}
}
