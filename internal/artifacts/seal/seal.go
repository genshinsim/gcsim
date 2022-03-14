package seal

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterSetFunc("seal of insulation", New)
	core.RegisterSetFunc("emblemofseveredfate", New)
	core.RegisterSetFunc("esr", New)
}

func New(c coretype.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.ER] = 0.20
		c.AddMod(coretype.CharStatMod{
			Key: "esr-2pc",
			Amount: func() ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		m := make([]float64, core.EndStatType)
		er := c.Stat(core.ER) + 1
		amt := 0.25 * er
		if amt > 0.75 {
			amt = 0.75
		}
		m[core.DmgP] = amt
		c.AddPreDamageMod(coretype.PreDamageMod{
			Key: "esr-4pc",
			Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
				if atk.Info.AttackTag == core.AttackTagElementalBurst {
					//calc er
					er := c.Stat(core.ER) + 1
					amt := 0.25 * er
					if amt > 0.75 {
						amt = 0.75
					}
					m[core.DmgP] = amt
					return m, true
				}
				return nil, false
			},
			Expiry: -1,
		})
	}
	//add flat stat to char
}
