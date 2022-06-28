package kuki

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) a1() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.Heal] = .15
	c.AddStatMod(character.StatMod{Base: modifier.NewBase("kuki-a1", -1), AffectedStat: attributes.Heal, Amount: func() ([]float64, bool) {
		if c.HPCurrent/c.MaxHP() <= 0.5 {
			return m, true
		}
		return nil, false
	}})
}
