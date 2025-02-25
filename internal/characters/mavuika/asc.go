package mavuika

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a1Key = "mavuika-a1"
	a4Key = "mavuika-a4"
)

func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.3
	c.Core.Events.Subscribe(event.OnNightsoulBurst, func(args ...interface{}) bool {
		c.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag(a1Key, 10*60),
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
		return false
	}, a1Key)
}

func (c *char) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}
	c.a4buff = make([]float64, attributes.EndStatType)
}

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	started := c.Core.F
	for _, char := range c.Core.Player.Chars() {
		this := char
		this.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag(a4Key, 20*60),
			Amount: func(_ *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
				// char must be active
				if c.Core.Player.Active() != this.Index {
					return nil, false
				}
				// Check hitlag extension to scale decay
				extension := c.StatusExpiry(a4Key) - (started + 20*60)
				dmg := c.burstStacks*0.002 + c.c4BonusVal()
				dmg *= 1.0 - float64(c.Core.F-started-extension)*c.c4DecayRate()
				c.a4buff[attributes.DmgP] = dmg
				return c.a4buff, true
			},
		})
	}
}
