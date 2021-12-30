package catch

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("the catch", weapon)
	core.RegisterWeaponFunc("thecatch", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {

	val := make([]float64, core.EndStatType)
	val[core.DmgP] = 0.12 + 0.04*float64(r)
	val[core.CR] = 0.045 + 0.015*float64(r)

	char.AddMod(core.CharStatMod{
		Key:    "the-catch",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if a == core.AttackTagElementalBurst {
				return val, true
			}
			return nil, false
		},
	})
	return "thecatch"
}
