package homa

import (
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("staff of homa", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	//add on hit effect to sim?
	m := make([]float64, core.EndStatType)
	m[core.HPP] = 0.15 + float64(r)*0.05
	atkp := 0.006 + float64(r)*0.002
	lowhp := 0.008 + float64(r)*0.002

	char.AddMod(core.CharStatMod{
		Key: "homa hp bonus",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			per := atkp
			if char.HP()/char.MaxHP() <= 0.5 {
				per += lowhp
			}
			m[core.ATKP] = per * char.MaxHP()
			return m, true
		},
		Expiry: -1,
	})

}
