package xiao

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) a4() {
	m := make([]float64, attributes.EndStatType)
	c.AddStatMod(character.StatMod{Base: modifier.NewBase("xiao-a4", -1), AffectedStat: attributes.DmgP, Amount: func() ([]float64, bool) {
		stacks := c.Tags["a4"]
		if stacks == 0 {
			return nil, false
		}
		m[attributes.DmgP] = float64(stacks) * 0.15
		return m, true
	}})
}
