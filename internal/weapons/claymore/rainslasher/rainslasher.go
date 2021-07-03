package rainslasher

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("rainslasher", weapon)
}

//Increases DMG against enemies affected by Hydro or Electro by 20/24/28/32/36%.
func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	dmg := 0.16 + float64(r)*0.04

	s.AddOnAttackWillLand(func(t def.Target, ds *def.Snapshot) {
		switch t.AuraType() {
		case def.Hydro:
		case def.Electro:
		case def.EC:
		default:
			return
		}
		ds.Stats[def.DmgP] += dmg
		log.Debugw("rainslasher", "frame", s.Frame(), "event", def.LogCalc, "final dmg%", ds.Stats[def.DmgP])
	}, fmt.Sprintf("rainslasher-%v", c.Name()))

}
