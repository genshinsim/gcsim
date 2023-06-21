package bennett

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c2() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.ER] = .3

	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("bennett-c2", -1),
		AffectedStat: attributes.ER,
		Amount: func() ([]float64, bool) {
			return m, c.CurrentHPRatio() < 0.7
		},
	})
}
