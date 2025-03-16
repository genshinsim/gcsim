package mavuika

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a1Key    = "mavuika-a1"
	a4Key    = "mavuika-a4"
	a4BufKey = "mavuika-a4-buff"
	a4Dur    = 20 * 60
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
	for _, char := range c.Core.Player.Chars() {
		// make sure variable isn't mutated by later loops
		char := char
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase(a4BufKey, -1),
			Amount: func(_ *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
				// char must be active
				if c.Core.Player.Active() != char.Index {
					return nil, false
				}
				if !c.StatusIsActive(a4Key) {
					return nil, false
				}
				dmg := c.burstStacks*0.002 + c.c4BonusVal()
				dmg *= float64(c.a4stacks) / 20.0
				c.a4buff[attributes.DmgP] = dmg
				return c.a4buff, true
			},
		})
	}
}

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	c.AddStatus(a4Key, a4Dur, true)
	c.a4stacks = 20
	c.a4src = c.Core.F
	c.a4DecayTask(c.a4src)
}

func (c *char) a4DecayTask(src int) {
	c.QueueCharTask(func() {
		if c.a4src != src {
			return
		}
		if !c.StatusIsActive(a4Key) {
			return
		}
		if c.a4stacks <= 0 {
			return
		}
		c.a4stacks -= c.c4DecayRate()
		c.a4DecayTask(src)
	}, 60)
}
