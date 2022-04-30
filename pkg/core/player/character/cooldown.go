package character

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

type CooldownModFunc func(a action.Action) float64

type cooldownMod struct {
	Amount func(a action.Action) float64
	modTmpl
}

func (c *CharWrapper) CDReduction(a action.Action, dur int) int {
	var cd float64 = 1
	n := 0
	for _, v := range c.cooldownMods {
		//if not expired
		if v.expiry == -1 || v.expiry > *c.f {
			amt := v.Amount(a)
			c.log.NewEvent(
				"applying cooldown modifier",
				glog.LogActionEvent,
				c.Index,
				"key", v.Key,
				"modifier", amt,
				"expiry", v.Expiry,
			)
			cd += amt
			c.cooldownMods[n] = v
			n++
		}
	}
	c.cooldownMods = c.cooldownMods[:n]

	return int(float64(dur) * cd)
}

func (c *CharWrapper) AddCooldownMod(key string, dur int, f CooldownModFunc) {
	expiry := *c.f + dur
	if dur < 0 {
		expiry = -1
	}
	mod := cooldownMod{
		modTmpl: modTmpl{
			key:    key,
			expiry: expiry,
		},
		Amount: f,
	}
	addMod(c, c.cooldownMods, &mod)
}

func (c *CharWrapper) DeleteCooldowntMod(key string) {
	deleteMod(c, c.cooldownMods, key)
}

func (c *CharWrapper) CooldownModIsActive(key string, char int, a action.Action) bool {
	_, ok := findModCheckExpiry(c.cooldownMods, key, *c.f)
	return ok
}
