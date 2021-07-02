package generic

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("rut", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	m := make([]float64, def.EndStatType)
	inc := .3 + float64(r)*0.1
	c.AddMod(def.CharStatMod{
		Key: "rust",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			if a == def.AttackTagNormal {
				m[def.DmgP] = inc
				return m, true
			}
			if a == def.AttackTagExtra {
				m[def.DmgP] = -0.1
				return m, true
			}
			return nil, false
		},
		Expiry: -1,
	})
}
