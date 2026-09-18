package iansan

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1ICDKey = "iansan-c1-icd"
	c2Key    = "iansan-c2"
	c6Key    = "iansan-c6"
)

func (c *char) c1(points float64) {
	if c.Base.Cons < 1 {
		return
	}
	if c.StatusIsActive(c1ICDKey) {
		return
	}
	c.c1Points += points
	if c.c1Points < 6 {
		return
	}

	c.AddEnergy("iansan-c1", 15)
	c.AddStatus(c1ICDKey, 18*60, true)
	c.c1Points = 0
}

func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}
	if c.Base.Ascension < 1 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.3

	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) {
		prev := args[0].(int)
		next := args[1].(int)

		buffDur := c.StatusDuration(a1Key)
		if buffDur == 0 {
			return
		}

		if prev != c.Index() {
			c.Core.Player.ByIndex(prev).DeleteStatMod(c2Key)
		}

		if next != c.Index() {
			// we don't need to worry about granting this buff when iansan gains her A1, because she
			// will always be the active character (using skill or using burst)
			c.Core.Player.ByIndex(next).AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(c2Key, buffDur),
				AffectedStat: attributes.ATKP,
				Amount: func() []float64 {
					if !c.StatModIsActive(a1Key) {
						return nil
					}
					return m
				},
			})
		}
	}, c2Key)
}

func (c *char) c2OnBurst() {
	if c.Base.Cons < 2 {
		return
	}

	c.a1AddBuff()
}

func (c *char) c4Init() {
	if c.Base.Cons < 4 {
		return
	}

	c.Core.Events.Subscribe(event.OnBurst, func(args ...interface{}) {
		if !c.StatusIsActive(burstStatus) {
			return
		}
		if c.Index() == c.Core.Player.Active() {
			return
		}
		if c.c4Generated {
			return
		}
		c.c4Stacks = 2
		c.c4Generated = true
	}, "iansan-c4")
}

func (c *char) c4OnBurst() {
	if c.Base.Cons < 4 {
		return
	}
	c.c4Generated = false
	c.c4Stacks = 0
}

func (c *char) c4Points() float64 {
	if c.Base.Cons < 4 {
		return 0.0
	}
	points := c.pointsOverflow * 0.5
	if c.c4Stacks > 0 {
		points += 4.0
		c.c4Stacks--
	}
	return points
}

func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}

	c.c6Buff = make([]float64, attributes.EndStatType)
	c.c6Buff[attributes.DmgP] = 0.25
}

func (c *char) c6BurstBonusDur() int {
	if c.Base.Cons < 6 {
		return 0
	}

	return 3 * 60
}

func (c *char) c6OnOverflow() {
	if c.Base.Cons < 6 {
		return
	}

	active := c.Core.Player.ActiveChar()
	active.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag(c6Key, 3*60),
		Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
			return c.c6Buff
		},
	})
}
