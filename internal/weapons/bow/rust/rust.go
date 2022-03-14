package generic

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("rust", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	m := make([]float64, core.EndStatType)
	inc := .3 + float64(r)*0.1
	char.AddPreDamageMod(coretype.PreDamageMod{
		Key: "rust",
		Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
			if atk.Info.AttackTag == coretype.AttackTagNormal {
				m[core.DmgP] = inc
				return m, true
			}
			if atk.Info.AttackTag == coretype.AttackTagExtra {
				m[core.DmgP] = -0.1
				return m, true
			}
			return nil, false
		},
		Expiry: -1,
	})

	return "rust"
}
