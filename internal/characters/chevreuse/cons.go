package mika

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.PyroP] = 0.20
	m[attributes.ElectroP] = 0.20

	buffDuration := 8
	active := c.Core.Player.ActiveChar()

	active.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("chev-c6", buffDuration),
		AffectedStat: attributes.NoStat,
		Extra:        true,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

}
