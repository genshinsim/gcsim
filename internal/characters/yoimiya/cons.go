package yoimiya

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

func (c *char) c1() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.2
	c.Core.Events.Subscribe(event.OnTargetDied, func(args ...interface{}) bool {
		//we assume target is affected if it's active
		if c.Core.Status.Duration("aurous") <= 0 {
			return false
		}

		c.AddStatMod("yoimiya-c1", 1200, attributes.ATKP, func() ([]float64, bool) {
			return m, true
		})

		return false
	}, "yoimiya-c1")
}

func (c *char) c2() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.PyroP] = 0.25
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		crit := args[3].(bool)

		if atk.Info.ActorIndex != c.Index || !crit {
			return false
		}

		c.AddStatMod("yoimiya-c2", 360, attributes.PyroP, func() ([]float64, bool) {
			return m, true
		})

		return false
	}, "yoimiya-c2")
}
