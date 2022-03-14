package ganyu

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func (c *char) c1() {
	c.Core.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*coretype.AttackEvent)
		t := args[0].(coretype.Target)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if atk.Info.Abil != "Frost Flake Arrow" {
			return false
		}
		c.AddEnergy("ganyu-c1", 2)
		t.AddResMod("ganyu-c1", core.ResistMod{
			Ele:      coretype.Cryo,
			Value:    -0.15,
			Duration: 5 * 60,
		})
		return false
	}, "ganyu-c1")

}
