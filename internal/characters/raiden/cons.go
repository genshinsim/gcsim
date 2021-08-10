package raiden

import "github.com/genshinsim/gsim/pkg/core"

func (c *char) c6() {
	c.Sim.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
		if c.Sim.Status("raidenburst") == 0 {
			return
		}
		if ds.ActorIndex != c.Index {
			return
		}
		if ds.Abil != "Musou Isshin" {
			return
		}
		if c.c6ICD > c.Sim.Frame() {
			return
		}
		if c.c6Count == 5 {
			return
		}
		c.c6ICD = c.Sim.Frame() + 60
		c.c6Count++
		for _, char := range c.Sim.Characters() {
			char.ReduceActionCooldown(core.ActionBurst, 1)
		}

	}, "raiden-c6")
}