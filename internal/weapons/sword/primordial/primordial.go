package primordial

import (
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("primordial jade cutter", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	//add on hit effect to sim?
	m := make([]float64, core.EndStatType)
	m[core.HPP] = 0.15 + float64(r)*0.05
	atkp := 0.009 + float64(r)*0.003

	char.AddMod(core.CharStatMod{
		Key: "homa hp bonus",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			m[core.ATKP] = atkp * char.MaxHP()
			return m, true
		},
		Expiry: -1,
	})

}
