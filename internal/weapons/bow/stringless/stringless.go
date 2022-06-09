package stringless

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("the stringless", weapon)
	core.RegisterWeaponFunc("thestringless", weapon)
	core.RegisterWeaponFunc("stringless", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	m := make([]float64, core.EndStatType)
	m[core.DmgP] = 0.18 + float64(r)*0.06
	char.AddPreDamageMod(core.PreDamageMod{
		Key: "stringless",
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			switch atk.Info.AttackTag {
			case core.AttackTagElementalArt:
			case core.AttackTagElementalArtHold:
			case core.AttackTagElementalBurst:
			default:
				return nil, false
			}
			return m, true
		},
		Expiry: -1,
	})
	return "thestringless"
}
