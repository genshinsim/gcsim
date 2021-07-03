package magicguide

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("magic guide", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	dmg := 0.09 + float64(r)*0.03

	s.AddOnAttackWillLand(func(t def.Target, ds *def.Snapshot) {
		if t.AuraContains(def.Hydro, def.Electro, def.Cryo) {
			ds.Stats[def.DmgP] += dmg
			log.Debugw("magic guide", "frame", s.Frame(), "event", def.LogCalc, "final dmg%", ds.Stats[def.DmgP])
		}
	}, fmt.Sprintf("magic-guide-%v", c.Name()))
}
