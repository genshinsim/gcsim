package magicguide

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("magic guide", weapon)
	core.RegisterWeaponFunc("magicguide", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	dmg := 0.09 + float64(r)*0.03

	c.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		t := args[0].(core.Target)
		ds := args[1].(*core.Snapshot)

		if t.AuraContains(core.Hydro, core.Electro, core.Cryo) {
			ds.Stats[core.DmgP] += dmg
			c.Log.Debugw("magic guide", "frame", c.F, "event", core.LogCalc, "final dmg%", ds.Stats[core.DmgP])
		}
		return false
	}, fmt.Sprintf("magic-guide-%v", char.Name()))
}
