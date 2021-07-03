package festering

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("festering desire", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {

	m := make([]float64, def.EndStatType)
	m[def.CR] = .045 + .015*float64(r)
	m[def.DmgP] = .12 + 0.04*float64(r)
	c.AddMod(def.CharStatMod{
		Key:    "festering",
		Expiry: -1,
		Amount: func(a def.AttackTag) ([]float64, bool) {
			return m, a == def.AttackTagElementalArt
		},
	})

}
