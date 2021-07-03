package homa

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("staff of homa", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	//add on hit effect to sim?
	m := make([]float64, def.EndStatType)
	m[def.HPP] = 0.15 + float64(r)*0.05
	atkp := 0.006 + float64(r)*0.002
	lowhp := 0.008 + float64(r)*0.002

	c.AddMod(def.CharStatMod{
		Key: "homa hp bonus",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			per := atkp
			if c.HP()/c.MaxHP() <= 0.5 {
				per += lowhp
			}
			m[def.ATKP] = per * c.MaxHP()
			return m, true
		},
		Expiry: -1,
	})

}
