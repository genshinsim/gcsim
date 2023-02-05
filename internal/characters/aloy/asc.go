package aloy

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// When Aloy receives the Coil effect from Frozen Wilds, her ATK is increased by 16%, while nearby party members' ATK is increased by 8%.
// This effect lasts 10s.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	for _, char := range c.Core.Player.Chars() {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ATKP] = .08
		if char.Index == c.Index {
			m[attributes.ATKP] = .16
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("aloy-a1", rushingIceDuration),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
}

// When Aloy is in the Rushing Ice state conferred by Frozen Wilds, her Cryo DMG Bonus increases by 3.5% every 1s.
// A maximum Cryo DMG Bonus increase of 35% can be gained in this way.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	stacks := 1
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("aloy-strong-strike", rushingIceDuration),
		AffectedStat: attributes.CryoP,
		Amount: func() ([]float64, bool) {
			if stacks > 10 {
				stacks = 10
			}
			m[attributes.CryoP] = float64(stacks) * 0.035
			return m, true
		},
	})

	for i := 0; i < 10; i++ {
		// every 1s, aloy can't experience hitlag so this way is fine
		c.Core.Tasks.Add(func() { stacks++ }, 60*(1+i))
	}
}
