package klee

import (
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (c *char) a1(a combat.AttackCB) {
	if c.Core.F < c.sparkICD {
		return
	}
	if c.Core.Rand.Float64() < 0.5 {
		return
	}
	c.sparkICD = c.Core.F + 60*4
	c.Core.Status.Add("kleespark", 60*30)
	c.Core.Log.NewEvent("klee gained spark", glog.LogCharacterEvent, c.Index).
		Write("icd", c.sparkICD)
}

func (c *char) a4() {
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		crit := args[3].(bool)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagExtra {
			return false
		}
		if !crit {
			return false
		}
		for _, x := range c.Core.Player.Chars() {
			x.AddEnergy("klee-a4", 2)
		}
		return false
	}, "kleea1")
}
