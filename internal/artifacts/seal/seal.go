package seal

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterSetFunc("seal of insulation", New)
}

func New(c def.Character, s def.Sim, log def.Logger, count int) {
	if count >= 2 {
		m := make([]float64, def.EndStatType)
		m[def.ER] = 0.20
		c.AddMod(def.CharStatMod{
			Key: "seal-2pc",
			Amount: func(a def.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		m := make([]float64, def.EndStatType)
		c.AddMod(def.CharStatMod{
			Key: "seal-4pc",
			Amount: func(ds def.AttackTag) ([]float64, bool) {
				if ds == def.AttackTagElementalBurst {
					//calc er
					er := c.Stat(def.ER) + 1
					amt := 0.3 * er
					if amt > 0.75 {
						amt = 0.75
					}
					m[def.DmgP] = amt
					return m, true
				}
				return nil, false
			},
			Expiry: s.Frame() + 600,
		})
	}
	//add flat stat to char
}
