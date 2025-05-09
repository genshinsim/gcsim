package iansan

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1ICD    = "iansan-c1"
	c6Status = "iansan-c6"
)

func (c *char) c1(points float64) {
	if c.Base.Cons < 1 {
		return
	}
	if c.StatusIsActive(c1ICD) {
		return
	}
	c.c1Points += points
	if c.c1Points < 6 {
		return
	}
	c.AddEnergy("iansan-c1", 15)
	c.AddStatus(c1ICD, 18*60, true)
}

func (c *char) c2ATKBuff(char *character.CharWrapper) {
	c.burstBuff[attributes.ATKP] = 0
	if c.Base.Cons < 2 {
		return
	}
	if c.Base.Ascension < 1 {
		return
	}
	if c.Index == c.Core.Player.Active() {
		return
	}
	if char.Index != c.Core.Player.Active() {
		return
	}
	c.burstBuff[attributes.ATKP] = 0.3
}

func (c *char) c4() {
	if c.Base.Cons < 4 {
		return
	}

	c.Core.Events.Subscribe(event.OnBurst, func(args ...interface{}) bool {
		if !c.StatusIsActive(burstStatus) {
			return false
		}
		if c.Index == c.Core.Player.Active() {
			return false
		}
		c.c4Stacks = 2
		return false
	}, "iansan-c4")
}

func (c *char) c4Points() float64 {
	if c.Base.Cons < 4 {
		return 0.0
	}
	points := c.pointsOverflow * 0.5
	if c.c4Stacks > 0 {
		points += 4.0
	}
	return points
}

func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.25

	active := c.Core.Player.ActiveChar()
	active.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag(c6Status, 3*60),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			return m, true
		},
	})
}
