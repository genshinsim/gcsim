package dragonbane

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("lion's roar", weapon)
}

//Increases DMG against enemies affected by Hydro or Electro by 20/24/28/32/36%.
func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	dmg := 0.16 + float64(r)*0.04

	s.AddOnAttackWillLand(func(t def.Target, ds *def.Snapshot) {
		if t.AuraContains(def.Electro, def.Pyro) {
			ds.Stats[def.DmgP] += dmg
			log.Debugw("lion's roar", "frame", s.Frame(), "event", def.LogCalc, "final dmg%", ds.Stats[def.DmgP])
		}
	}, fmt.Sprintf("lion-%v", c.Name()))

}
