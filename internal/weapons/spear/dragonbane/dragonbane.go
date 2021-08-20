package dragonbane

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("dragon's bane", weapon)
}

//Increases DMG against enemies affected by Hydro or Electro by 20/24/28/32/36%.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	dmg := 0.16 + float64(r)*0.04

	c.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		t := args[0].(core.Target)
		if ds.ActorIndex != char.CharIndex() {
			return false
		}
		// if t.AuraType() == def.Hydro {
		// 	ds.Stats[def.DmgP] += dmg
		// 	c.Log.Debugw("dragonbane", "frame", c.F, "event", def.LogCalc, "final dmg%", ds.Stats[def.DmgP])
		// }
		if t.AuraContains(core.Hydro, core.Pyro) {
			ds.Stats[core.DmgP] += dmg
			c.Log.Debugw("dragonbane", "frame", c.F, "event", core.LogCalc, "final dmg%", ds.Stats[core.DmgP])
		}
		return false
	}, fmt.Sprintf("dragonbane-%v", char.Name()))

}
