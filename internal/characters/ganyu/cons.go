package ganyu

import (
	"github.com/genshinsim/gsim/pkg/def"
)

func (c *char) c1() {
	c.Sim.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.Index {
			return
		}
		if ds.Abil != "Frost Flake Arrow" {
			return
		}
		c.AddEnergy(2)
		t.AddResMod("ganyu-c1", def.ResistMod{
			Ele:      def.Cryo,
			Value:    -0.15,
			Duration: 5 * 60,
		})

	}, "ganyu-c1")
}
