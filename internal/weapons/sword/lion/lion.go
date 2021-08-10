package dragonbane

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("lion's roar", weapon)
}

//Increases DMG against enemies affected by Hydro or Electro by 20/24/28/32/36%.
func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	dmg := 0.16 + float64(r)*0.04

	s.AddOnAttackWillLand(func(t core.Target, ds *core.Snapshot) {
		if t.AuraContains(core.Electro, core.Pyro) {
			ds.Stats[core.DmgP] += dmg
			log.Debugw("lion's roar", "frame", s.Frame(), "event", core.LogCalc, "final dmg%", ds.Stats[core.DmgP])
		}
	}, fmt.Sprintf("lion-%v", c.Name()))

}
