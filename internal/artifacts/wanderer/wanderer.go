package wanderer

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterSetFunc("wanderer's troupe", New)
}

func New(c core.Character, s core.Sim, log core.Logger, count int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.EM] = 80
		c.AddMod(core.CharStatMod{
			Key: "wt-2pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		switch c.WeaponClass() {
		case core.WeaponClassCatalyst:
		case core.WeaponClassBow:
		default:
			//don't add this mod if wrong weapon class
			return
		}
		m := make([]float64, core.EndStatType)
		m[core.DmgP] = 0.35
		c.AddMod(core.CharStatMod{
			Key: "wt-4pc",
			Amount: func(ds core.AttackTag) ([]float64, bool) {
				if ds != core.AttackTagNormal && ds != core.AttackTagExtra {
					return nil, false
				}
				return m, true
			},
			Expiry: -1,
		})
	}
	//add flat stat to char
}
