package xiangling

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

func (c *char) c1(a combat.AttackCB) {
	if c.Base.Cons < 1 {
		return
	}
	e, ok := a.Target.(core.Enemy)
	if !ok {
		return
	}
	e.AddResistMod("xiangling-c1", 6*60, attributes.Pyro, -0.15)
}

func (c *char) c6(dur int) {
	m := make([]float64, attributes.EndStatType)
	m[attributes.PyroP] = 0.15

	c.Core.Status.Add("xlc6", dur)

	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod("xiangling-c6", dur, attributes.PyroP, func() ([]float64, bool) {
			return m, true
		})
	}
}
