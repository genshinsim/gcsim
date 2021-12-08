package raiden

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) c6() {
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if c.Core.Status.Duration("raidenburst") == 0 {
			return false
		}
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if atk.Info.Abil != "Musou Isshin" {
			return false
		}
		if c.c6ICD > c.Core.F {
			return false
		}
		if c.c6Count == 5 {
			return false
		}
		c.c6ICD = c.Core.F + 60
		c.c6Count++
		for _, char := range c.Core.Chars {
			char.ReduceActionCooldown(core.ActionBurst, 1)
		}
		return false
	}, "raiden-c6")
}
