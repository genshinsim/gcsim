package dragonbane

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("dragon's bane", weapon)
}

//Increases DMG against enemies affected by Hydro or Electro by 20/24/28/32/36%.
func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	dmg := 0.16 + float64(r)*0.04

	s.AddOnAttackWillLand(func(t core.Target, ds *core.Snapshot) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		// if t.AuraType() == def.Hydro {
		// 	ds.Stats[def.DmgP] += dmg
		// 	log.Debugw("dragonbane", "frame", s.Frame(), "event", def.LogCalc, "final dmg%", ds.Stats[def.DmgP])
		// }
		if t.AuraContains(core.Hydro, core.Pyro) {
			ds.Stats[core.DmgP] += dmg
			log.Debugw("dragonbane", "frame", s.Frame(), "event", core.LogCalc, "final dmg%", ds.Stats[core.DmgP])
		}
	}, fmt.Sprintf("dragonbane-%v", c.Name()))

}
