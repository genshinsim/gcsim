package xingqiu

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) a4() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.HydroP] = 0.2

	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("xingqiu-a4", -1),
		AffectedStat: attributes.HydroP,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}
