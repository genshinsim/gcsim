package ganyu

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) c1() {
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		t := args[0].(core.Target)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if atk.Info.Abil != "Frost Flake Arrow" {
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
