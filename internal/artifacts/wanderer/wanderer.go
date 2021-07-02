package wanderer

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterSetFunc("wanderer's troupe", New)
}

func New(c def.Character, s def.Sim, log def.Logger, count int) {
	if count >= 2 {
		m := make([]float64, def.EndStatType)
		m[def.EM] = 80
		c.AddMod(def.CharStatMod{
			Key: "wt-2pc",
			Amount: func(a def.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		switch c.WeaponClass() {
		case def.WeaponClassCatalyst:
		case def.WeaponClassBow:
		default:
			//don't add this mod if wrong weapon class
			return
		}
		m := make([]float64, def.EndStatType)
		m[def.DmgP] = 0.35
		c.AddMod(def.CharStatMod{
			Key: "glad-4pc",
			Amount: func(ds def.AttackTag) ([]float64, bool) {
				if ds != def.AttackTagNormal && ds != def.AttackTagExtra {
					return nil, false
				}
				return m, true
			},
			Expiry: -1,
		})
	}
	//add flat stat to char
}
