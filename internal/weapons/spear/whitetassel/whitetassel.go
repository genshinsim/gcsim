package whitetassel

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("whitetassel", weapon)
}

// Increases Normal Attack DMG by 24/30/36/42/48%.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	m := make([]float64, core.EndStatType)
	m[core.DmgP] = 0.18 + 0.06*float64(r)

	char.AddPreDamageMod(core.PreDamageMod{
		Key:    "whitetassel",
		Expiry: -1,
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			return m, atk.Info.AttackTag == core.AttackTagNormal
		},
	})

	return "whitetassel"
}
