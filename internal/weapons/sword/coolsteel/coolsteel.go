package coolsteel

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("coolsteel", weapon)
}

//Increases DMG against enemies affected by Hydro or Cryo by 12-24%.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	dmg := 0.09 + float64(r)*0.03

	char.AddPreDamageMod(core.PreDamageMod{
		Key:    "coolsteel",
		Expiry: -1,
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			m := make([]float64, core.EndStatType)
			if t.AuraContains(core.Hydro, core.Cryo) {
				m[core.DmgP] = dmg
				return m, true
			}
			return nil, false
		},
	})
	return "coolsteel"
}
