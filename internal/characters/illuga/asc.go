package illuga

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var (
	a4GeoBuff = []float64{0, 0.07, 0.14, 0.24}
	a4LcrBuff = []float64{0, 0.48, 0.96, 1.6}
)

func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = 0.05 + c.c6CR()
	m[attributes.CD] = 0.1 + c.c6CD()

	if c.Core.Player.GetMoonsignLevel() >= 2 {
		m[attributes.EM] = 50 + c.c6EM()
	}

	for _, char := range c.Core.Player.Chars() {
		if char.Index() == c.Index() {
			continue
		}
		char.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag("illuga-a1", 20*60),
			Amount: func() []float64 {
				return m
			},
		})
	}
}

func (c *char) a4Count() int {
	if c.Base.Ascension < 4 {
		return 0
	}

	result := 0
	for _, char := range c.Core.Player.Chars() {
		switch char.Base.Element {
		case attributes.Geo, attributes.Hydro:
			result++
		}
	}
	return min(result, 3)
}

func (c *char) a4GeoBonus() float64 {
	return a4GeoBuff[c.a4Count()]
}

func (c *char) a4LcrBonus() float64 {
	return a4LcrBuff[c.a4Count()]
}
