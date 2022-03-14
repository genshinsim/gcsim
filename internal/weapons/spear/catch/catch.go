package catch

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("the catch", weapon)
	core.RegisterWeaponFunc("thecatch", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {

	val := make([]float64, core.EndStatType)
	val[core.DmgP] = 0.12 + 0.04*float64(r)
	val[core.CR] = 0.045 + 0.015*float64(r)

	char.AddPreDamageMod(coretype.PreDamageMod{
		Key:    "the-catch",
		Expiry: -1,
		Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
			if atk.Info.AttackTag == core.AttackTagElementalBurst {
				return val, true
			}
			return nil, false
		},
	})
	return "thecatch"
}
