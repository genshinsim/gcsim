package mualani

import (
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

const a1Delay = 20

func (c *char) a1cb() info.AttackCBFunc {
	if c.Base.Ascension < 1 {
		return nil
	}
	done := false
	return func(a info.AttackCB) {
		if c.a1Count >= 2 {
			return
		}
		if done {
			return
		}
		if a.Target.Type() != info.TargettableEnemy {
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
	c.Core.Events.Subscribe(event.OnNightsoulBurst, func(args ...any) bool {
		c.a4Stacks = min(c.a4Stacks+1, 3)
		return false
	}, "maulani-a4-nightsoul")
}
