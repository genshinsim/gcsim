package yoimiya

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) c1() {
	m := make([]float64, core.EndStatType)
	m[core.ATKP] = 0.2
	c.Core.Events.Subscribe(core.OnTargetDied, func(args ...interface{}) bool {
		//we assume target is affected if it's active
		if c.Core.Status.Duration("aurous") > 0 {
			c.AddMod(core.CharStatMod{
				Key:    "yoimiya-c1",
				Expiry: c.Core.F + 1200,
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})
		}
		return false
	}, "yoimiya-c1")
}

func (c *char) c2() {
	m := make([]float64, core.EndStatType)
	m[core.PyroP] = 0.25
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		crit := args[3].(bool)
		if atk.Info.ActorIndex == c.Index && crit {
			c.AddMod(core.CharStatMod{
				Key:    "yoimiya-c2",
				Expiry: c.Core.F + 360,
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})
		}
		return false
	}, "yoimiya-c2")
}
