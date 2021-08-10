package magicguide

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("magic guide", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	dmg := 0.09 + float64(r)*0.03

	s.AddOnAttackWillLand(func(t core.Target, ds *core.Snapshot) {
		if t.AuraContains(core.Hydro, core.Electro, core.Cryo) {
			ds.Stats[core.DmgP] += dmg
			log.Debugw("magic guide", "frame", s.Frame(), "event", core.LogCalc, "final dmg%", ds.Stats[core.DmgP])
		}
	}, fmt.Sprintf("magic-guide-%v", c.Name()))
}
