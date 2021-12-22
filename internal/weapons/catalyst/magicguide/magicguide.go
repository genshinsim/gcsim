package magicguide

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("magic guide", weapon)
	core.RegisterWeaponFunc("magicguide", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	// dmg := 0.09 + float64(r)*0.03

	char.AddPreDamageMod(core.PreDamageMod{
		Key:    "magic-guide",
		Expiry: -1,
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			m := make([]float64, core.EndStatType)
			if t.AuraContains(core.Hydro, core.Electro, core.Cryo) {
				m[core.DmgP] = 0.09 + float64(r)*0.03
				return m, true
			}
			return nil, false
		},
	})
}
