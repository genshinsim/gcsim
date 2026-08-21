package kachina

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	playercharacter "github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const a1Key = "kachina-a1-geo-buff"

func (c *char) initA1() {
	if c.Base.Ascension < 1 {
		return
	}

	bonus := make([]float64, attributes.EndStatType)
	bonus[attributes.DmgP] = 0.20
	for _, target := range c.Core.Player.Chars() {
		char := target
		char.AddAttackMod(playercharacter.AttackMod{
			Base: modifier.NewBase(a1Key+"-"+char.Base.Key.String(), -1),
			Amount: func(atk *info.AttackEvent, _ info.Target) []float64 {
				if atk.Info.Element != attributes.Geo || !char.StatusIsActive(a1Key) {
					return nil
				}
				return bonus
			},
		})
	}

	c.Core.Events.Subscribe(event.OnNightsoulBurst, func(_ ...any) {
		c.Core.Player.ActiveChar().AddStatus(a1Key, 12*60, true)
	}, a1Key+"-event")
}
