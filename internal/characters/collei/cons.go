package collei

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c1() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.ER] = 0.2
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("collei-c1", -1),
		AffectedStat: attributes.ER,
		Amount: func() ([]float64, bool) {
			if c.Core.Player.Active() != c.Index {
				return m, true
			}
			return nil, false
		},
	})
}
