package homa

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("staff of homa", weapon)
	core.RegisterWeaponFunc("staffofhoma", weapon)
	core.RegisterWeaponFunc("homa", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {

	mHP := make([]float64, core.EndStatType)
	mHP[core.HPP] = 0.15 + float64(r)*0.05
	char.AddMod(core.CharStatMod{
		Key:    "homa-hp",
		Expiry: -1,
		Amount: func() ([]float64, bool) {
			return mHP, true
		},
	})

	mATK := make([]float64, core.EndStatType)
	atkp := 0.006 + float64(r)*0.002
	lowhp := 0.008 + float64(r)*0.002
	char.AddMod(core.CharStatMod{
		Key:          "homa-atk-buff",
		Expiry:       -1,
		AffectedStat: core.ATK, // to avoid infinite loop when calling MaxHP
		Amount: func() ([]float64, bool) {
			maxhp := char.MaxHP()
			per := atkp
			if maxhp <= 0.5 {
				per += lowhp
			}
			mATK[core.ATK] = per * maxhp
			return mATK, true
		},
	})

	return "staffofhoma"
}
