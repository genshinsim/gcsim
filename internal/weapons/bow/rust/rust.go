package generic

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("rust", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	m := make([]float64, core.EndStatType)
	inc := .3 + float64(r)*0.1
	char.AddPreDamageMod(core.PreDamageMod{
		Key: "rust",
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			if atk.Info.AttackTag == core.AttackTagNormal {
				m[core.DmgP] = inc
				return m, true
			}
			if atk.Info.AttackTag == core.AttackTagExtra {
				m[core.DmgP] = -0.1
				return m, true
			}
			return nil, false
		},
		Expiry: -1,
	})

	return "rust"
}
