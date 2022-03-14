package whitetassel

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("whitetassel", weapon)
}

// Increases Normal Attack DMG by 24/30/36/42/48%.
func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	m := make([]float64, core.EndStatType)
	m[core.DmgP] = 0.18 + 0.06*float64(r)

	char.AddPreDamageMod(coretype.PreDamageMod{
		Key:    "whitetassel",
		Expiry: -1,
		Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
			return m, atk.Info.AttackTag == coretype.AttackTagNormal
		},
	})

	return "whitetassel"
}
