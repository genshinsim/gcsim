package stringless

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("the stringless", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	m := make([]float64, core.EndStatType)
	m[core.DmgP] = 0.18 + float64(r)*0.06
	c.AddMod(core.CharStatMod{
		Key: "stringless",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			switch a {
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
}
