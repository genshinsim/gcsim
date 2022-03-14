package wanderer

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterSetFunc("wanderer's troupe", New)
	core.RegisterSetFunc("wandererstroupe", New)
}

func New(c coretype.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.EM] = 80
		c.AddMod(coretype.CharStatMod{
			Key: "wt-2pc",
			Amount: func() ([]float64, bool) {
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
		c.AddPreDamageMod(coretype.PreDamageMod{
			Key: "wt-4pc",
			Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
				return m, atk.Info.AttackTag == coretype.AttackTagNormal || atk.Info.AttackTag == coretype.AttackTagExtra
			},
			Expiry: -1,
		})
	}
	//add flat stat to char
}
