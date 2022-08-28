package collei

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c4BuffKey = "collei-c4"

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

func (c *char) c2() {
	for _, event := range dendroEvents {
		c.Core.Events.Subscribe(event, func(args ...interface{}) bool {
			if c.c2Extended {
				return false
			}
			c.c2Extended = c.StatusIsActive(sproutKey) || c.StatusIsActive(skillKey)
			if c.StatusIsActive(sproutKey) {
				c.ExtendStatus(sproutKey, 180)
			}
			return false
		}, "collei-c2")
	}
}

func (c *char) c4() {
	for i, char := range c.Core.Player.Chars() {
		// does not affect collei
		if c.Index == i {
			continue
		}
		amts := make([]float64, attributes.EndStatType)
		amts[attributes.EM] = 80
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(c4BuffKey, 720),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return amts, true
			},
		})
	}
}
