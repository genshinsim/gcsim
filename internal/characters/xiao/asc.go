package xiao

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

func (c *char) a4() {
	m := make([]float64, attributes.EndStatType)
	c.AddStatMod("xiao-a4", -1, attributes.DmgP, func() ([]float64, bool) {
		stacks := c.Tags["a4"]
		if stacks == 0 {
			return nil, false
		}
		m[attributes.DmgP] = float64(stacks) * 0.15
		return m, true
	})
}
