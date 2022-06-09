package rainslasher

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("rainslasher", weapon)
}

//Increases DMG against enemies affected by Hydro or Electro by 20/24/28/32/36%.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	dmg := 0.16 + float64(r)*0.04

	char.AddPreDamageMod(core.PreDamageMod{
		Key:    "rainslasher",
		Expiry: -1,
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			m := make([]float64, core.EndStatType)
			if !t.AuraContains(core.Hydro, core.Electro) {
				return nil, false
			}
			m[core.DmgP] = dmg
			return m, true
		},
	})
	return "rainslasher"
}
