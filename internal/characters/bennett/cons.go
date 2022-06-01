package bennett

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

func (c *char) c2() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.ER] = .3

	c.AddStatMod("bennett-c2", -1, attributes.ER, func() ([]float64, bool) {
		return m, c.HPCurrent/c.MaxHP() < 0.7
	})
}
