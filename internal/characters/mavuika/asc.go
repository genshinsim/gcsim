package mavuika

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	c.Core.Events.Subscribe(event.OnNightsoulBurst, func(args ...interface{}) bool {
		buff := make([]float64, attributes.EndStatType)
		buff[attributes.ATKP] = 0.35
		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag("mavuika-a1", 10*60),
			Amount: func(a *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
				return buff, true
			},
		})
		return false
	}, "mavuika-a1-on-ns-burst")
}

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	started := c.Core.F
	c.baseA4Buff = min(0.5, 0.0025*c.consumedFightingSpirit)
	for _, char := range c.Core.Player.Chars() {
		this := char
		this.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("mavuika-a4", 20*60),
			Amount: func(_ *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
				// char must be active
				if c.Core.Player.Active() != this.Index {
					return nil, false
				}
				// buff time decay
				decay := 0.
				if c.Base.Cons < 4 {
					decay = float64((c.Core.F-started)/60) * (c.baseA4Buff) / 20
				}
				dmg := max(0., c.baseA4Buff-float64(decay))
				c.a4Buff[attributes.DmgP] = dmg
				return c.a4Buff, true
			},
		})
	}
}
