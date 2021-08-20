package ganyu

import (
	"github.com/genshinsim/gsim/pkg/core"
)

func (c *char) c1() {
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		t := args[0].(core.Target)
		if ds.ActorIndex != c.Index {
			return false
		}
		if ds.Abil != "Frost Flake Arrow" {
			return false
		}
		c.AddEnergy(2)
		t.AddResMod("ganyu-c1", core.ResistMod{
			Ele:      core.Cryo,
			Value:    -0.15,
			Duration: 5 * 60,
		})
		return false
	}, "ganyu-c1")

}
