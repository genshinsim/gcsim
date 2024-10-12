package mualani

import (
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

const a1Delay = 20

func (c *char) a1cb() combat.AttackCBFunc {
	if c.Base.Ascension < 1 {
		return nil
	}
	done := false
	return func(a combat.AttackCB) {
		if c.a1Count >= 2 {
			return
		}
		if done {
			return
		}
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		done = true
		c.a1Count++
		c.QueueCharTask(func() {
			if c.nightsoulState.HasBlessing() {
				c.nightsoulState.GeneratePoints(20)
				c.c2puffer()
				c.c4puffer()
			}
		}, a1Delay)
	}
}

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	c.Core.Events.Subscribe(event.OnNightsoulBurst, func(args ...interface{}) bool {
		c.a4Stacks = min(c.a4Stacks+1, 3)
		return false
	}, "maulani-a4")
}
func (c *char) a4amount() float64 {
	if c.Base.Ascension < 1 {
		return 0.0
	}
	s := c.a4Stacks
	c.a4Stacks = 0
	return 0.15 * float64(s) * c.MaxHP()
}
