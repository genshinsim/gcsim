package rainslasher

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("rainslasher", weapon)
}

//Increases DMG against enemies affected by Hydro or Electro by 20/24/28/32/36%.
func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	dmg := 0.16 + float64(r)*0.04

	s.AddOnAttackWillLand(func(t core.Target, ds *core.Snapshot) {
		if t.AuraContains(core.Hydro, core.Electro) {
			ds.Stats[core.DmgP] += dmg
			log.Debugw("rainslasher", "frame", s.Frame(), "event", core.LogCalc, "final dmg%", ds.Stats[core.DmgP])
		}
	}, fmt.Sprintf("rainslasher-%v", c.Name()))

}
