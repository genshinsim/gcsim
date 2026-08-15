package illuga

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var (
	a4GeoBuff = []float64{0, 0.07, 0.14, 0.24}
	a4LcrBuff = []float64{0, 0.48, 0.96, 1.6}
)

const (
	a1CR = 0.05
	a1CD = 0.1
	a1EM = 50
)

func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = a1CR + c.c6CR()
	m[attributes.CD] = a1CD + c.c6CD()

	n := make([]float64, attributes.EndStatType)

	if c.Core.Player.GetMoonsignLevel() >= 2 {
		n[attributes.EM] = a1EM + c.c6EM()
	}

	for _, char := range c.Core.Player.Chars() {
		if char.Index() == c.Index() {
			continue
		}
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag("illuga-a1-crit", 20*60),
			Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
				if atk.Info.Element != attributes.Geo {
					return nil
				}

				return m
			},
		})

		if c.Core.Player.GetMoonsignLevel() >= 2 {
			char.AddStatMod(character.StatMod{
				Base: modifier.NewBaseWithHitlag("illuga-a1-em", 20*60),
				Amount: func() []float64 {
					return n
				},
			})
		}
	}

	c.Core.Events.Subscribe(event.OnLunarReactionAttack, func(args ...any) {
		char := args[0].(*character.CharWrapper)
		if char.Index() == c.Index() {
			return
		}

		atk := args[1].(*info.AttackEvent)
		if atk.Info.AttackTag != attacks.AttackTagReactionLunarCrystallize {
			return
		}

		atk.Snapshot.Stats[attributes.CR] += a1CR + c.c6CR()
		atk.Snapshot.Stats[attributes.CD] += a1CD + c.c6CD()
	}, "illuga-a1-lunarcrystallize")
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
