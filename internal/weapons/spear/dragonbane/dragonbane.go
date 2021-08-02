package dragonbane

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("dragon's bane", weapon)
}

//Increases DMG against enemies affected by Hydro or Electro by 20/24/28/32/36%.
func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	dmg := 0.16 + float64(r)*0.04

	s.AddOnAttackWillLand(func(t def.Target, ds *def.Snapshot) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		// if t.AuraType() == def.Hydro {
		// 	ds.Stats[def.DmgP] += dmg
		// 	log.Debugw("dragonbane", "frame", s.Frame(), "event", def.LogCalc, "final dmg%", ds.Stats[def.DmgP])
		// }
		if t.AuraContains(def.Hydro, def.Pyro) {
			ds.Stats[def.DmgP] += dmg
			log.Debugw("dragonbane", "frame", s.Frame(), "event", def.LogCalc, "final dmg%", ds.Stats[def.DmgP])
		}
	}, fmt.Sprintf("dragonbane-%v", c.Name()))

}
