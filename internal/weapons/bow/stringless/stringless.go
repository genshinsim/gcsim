package stringless

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("the stringless", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	m := make([]float64, def.EndStatType)
	m[def.DmgP] = 0.18 + float64(r)*0.06
	c.AddMod(def.CharStatMod{
		Key: "stringless",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			switch a {
			case def.AttackTagElementalArt:
			case def.AttackTagElementalArtHold:
			case def.AttackTagElementalBurst:
			default:
				return nil, false
			}
			return m, true
		},
		Expiry: -1,
	})
}
