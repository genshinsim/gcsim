package kuki

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// When Shinobu's HP is not higher than 50%, her Healing Bonus is increased by 15%.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.Heal] = .15
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("kuki-a1", -1),
		AffectedStat: attributes.Heal,
		Amount: func() ([]float64, bool) {
			if c.HPCurrent/c.MaxHP() <= 0.5 {
				return m, true
			}
			return nil, false
		},
	})
}

// Sanctifying Ring's abilities will be boosted based on Shinobu's Elemental Mastery:
//
// - Healing amount will be increased by 75% of Elemental Mastery.
func (c *char) a4Healing() float64 {
	if c.Base.Ascension < 4 {
		return 0
	}
	return c.Stat(attributes.EM) * 0.75
}

// Sanctifying Ring's abilities will be boosted based on Shinobu's Elemental Mastery:
//
// - DMG dealt is increased by 25% of Elemental Mastery.
func (c *char) a4Damage() float64 {
	if c.Base.Ascension < 4 {
		return 0
	}
	return c.Stat(attributes.EM) * 0.25
}
