package dragonbane

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("dragon's bane", weapon)
	core.RegisterWeaponFunc("dragonsbane", weapon)
}

//Increases DMG against enemies affected by Hydro or Electro by 20/24/28/32/36%.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	dmg := 0.16 + float64(r)*0.04

	char.AddPreDamageMod(core.PreDamageMod{
		Key:    "dragonbane",
		Expiry: -1,
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			m := make([]float64, core.EndStatType)
			m[core.DmgP] = dmg
			if t.AuraContains(core.Hydro, core.Pyro) {
				return m, true
			}
			return nil, false
		},
	})

}
