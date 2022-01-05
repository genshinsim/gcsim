package gladiator

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("gladiator's finale", New)
	core.RegisterSetFunc("gladiatorsfinale", New)
	core.RegisterSetFunc("gladiators", New)
}

func New(c core.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.ATKP] = 0.18
		c.AddMod(core.CharStatMod{
			Key: "glad-2pc",
			Amount: func() ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		switch c.WeaponClass() {
		case core.WeaponClassSpear:
		case core.WeaponClassSword:
		case core.WeaponClassClaymore:
		default:
			//don't add this mod if wrong weapon class
			return
		}

		m := make([]float64, core.EndStatType)
		m[core.DmgP] = 0.35
		c.AddPreDamageMod(core.PreDamageMod{
			Key: "glad-4pc",
			Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
				return m, (atk.Info.AttackTag == core.AttackTagNormal || atk.Info.AttackTag == core.AttackTagExtra)
			},
			Expiry: -1,
		})
	}
	//add flat stat to char
}
