package ganyu

import (
	"github.com/genshinsim/gsim/pkg/core"
)

func (c *char) c1() {
	c.Sim.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.Index {
			return
		}
		if ds.Abil != "Frost Flake Arrow" {
			return
		}
		c.AddEnergy(2)
		t.AddResMod("ganyu-c1", core.ResistMod{
			Ele:      core.Cryo,
			Value:    -0.15,
			Duration: 5 * 60,
		})

	}, "ganyu-c1")
}
