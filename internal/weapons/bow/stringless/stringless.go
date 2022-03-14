package stringless

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("the stringless", weapon)
	core.RegisterWeaponFunc("thestringless", weapon)
	core.RegisterWeaponFunc("stringless", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	m := make([]float64, core.EndStatType)
	m[core.DmgP] = 0.18 + float64(r)*0.06
	char.AddPreDamageMod(coretype.PreDamageMod{
		Key: "stringless",
		Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
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
