package primordial

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("primordial jade cutter", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	//add on hit effect to sim?
	m := make([]float64, def.EndStatType)
	m[def.HPP] = 0.15 + float64(r)*0.05
	atkp := 0.009 + float64(r)*0.003

	c.AddMod(def.CharStatMod{
		Key: "homa hp bonus",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			m[def.ATKP] = atkp * c.MaxHP()
			return m, true
		},
		Expiry: -1,
	})

}
