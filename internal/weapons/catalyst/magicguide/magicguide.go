package magicguide

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("magic guide", weapon)
	core.RegisterWeaponFunc("magicguide", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	// dmg := 0.09 + float64(r)*0.03

	char.AddPreDamageMod(coretype.PreDamageMod{
		Key:    "magic-guide",
		Expiry: -1,
		Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
			m := make([]float64, core.EndStatType)
			if t.AuraContains(core.Hydro, core.Electro, coretype.Cryo) {
				m[core.DmgP] = 0.09 + float64(r)*0.03
				return m, true
			}
			return nil, false
		},
	})

	return "magicguide"
}
