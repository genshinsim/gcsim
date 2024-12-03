package mualani

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
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
	}, "maulani-a4-nightsoul")

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)

		switch atk.Info.AttackTag {
		case attacks.AttackTagElementalBurst:
		default:
			return false
		}
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		amt := 0.15 * float64(c.a4Stacks) * c.MaxHP()
		c.a4Stacks = 0

		c.Core.Log.NewEvent("mualani a4 proc dmg add", glog.LogPreDamageMod, atk.Info.ActorIndex).
			Write("before", atk.Info.FlatDmg).
			Write("addition", amt)

		atk.Info.FlatDmg += amt
		return false
	}, "maulani-a4-hook")
}
