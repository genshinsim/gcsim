package primordial

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("primordial jade cutter", weapon)
	core.RegisterWeaponFunc("primordialjadecutter", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	//add on hit effect to sim?
	m := make([]float64, core.EndStatType)
	m[core.HPP] = 0.15 + float64(r)*0.05
	atkp := 0.009 + float64(r)*0.003

	char.AddMod(core.CharStatMod{
		Key: "cutter hp bonus",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			m[core.ATK] = atkp * char.MaxHP()
			return m, true
		},
		Expiry: -1,
	})
	return "primordialjadecutter"
}
