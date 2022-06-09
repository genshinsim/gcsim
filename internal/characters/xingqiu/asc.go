package xingqiu

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

func (c *char) a4() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.HydroP] = 0.2

	c.AddStatMod("xingqiu-a4", -1, attributes.HydroP, func() ([]float64, bool) {
		return m, true
	})
}
