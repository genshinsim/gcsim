package kuki

import "github.com/genshinsim/gcsim/pkg/core/attributes"

func (c *char) a1() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.Heal] = .15
	c.AddStatMod("kuki-a1", -1, attributes.Heal, func() ([]float64, bool) {
		if c.HPCurrent/c.MaxHP() <= 0.5 {
			return m, true
		}
		return nil, false
	})
}
