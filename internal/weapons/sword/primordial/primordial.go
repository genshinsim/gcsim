package primordial

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("primordial jade cutter", weapon)
	core.RegisterWeaponFunc("primordialjadecutter", weapon)
	core.RegisterWeaponFunc("jadecutter", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {

	mHP := make([]float64, core.EndStatType)
	mHP[core.HPP] = 0.15 + float64(r)*0.05
	char.AddMod(core.CharStatMod{
		Key:    "jadecutter-hp",
		Expiry: -1,
		Amount: func() ([]float64, bool) {
			return mHP, true
		},
	})

	mATK := make([]float64, core.EndStatType)
	atkp := 0.009 + float64(r)*0.003
	char.AddMod(core.CharStatMod{
		Key:          "jadecutter-atk-buff",
		Expiry:       -1,
		AffectedStat: core.ATK, // to avoid infinite loop when calling MaxHP
		Amount: func() ([]float64, bool) {
			mATK[core.ATK] = atkp * char.MaxHP()
			return mATK, true
		},
	})

	return "primordialjadecutter"
}
