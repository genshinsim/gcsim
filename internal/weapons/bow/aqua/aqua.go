package aqua

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("aqua simulacra", weapon)
	core.RegisterWeaponFunc("aquasimulacra", weapon)
	core.RegisterWeaponFunc("aqua", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	//add on hit effect to sim?
	m := make([]float64, core.EndStatType)
	v := make([]float64, core.EndStatType)
	v[core.HPP] = 0.12 + float64(r)*0.04
	m[core.DmgP] = 0.15 + float64(r)*0.05

	char.AddMod(core.CharStatMod{
		Key: "aquasimulacra",
		Amount: func() ([]float64, bool) {
			return v, true
		},
		Expiry: -1,
	})

	//TODO: need range check here
	char.AddPreDamageMod(core.PreDamageMod{
		Key: "aquasimulacra",
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})
	return "aquasimulacra"
}
