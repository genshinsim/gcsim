package catch

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("the catch", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {

	val := make([]float64, core.EndStatType)
	val[core.DmgP] = 0.12 + 0.04*float64(r)
	val[core.CR] = 0.045 + 0.015*float64(r)

	c.AddMod(core.CharStatMod{
		Key: "the-catch",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if a == core.AttackTagElementalBurst {
				return val, true
			}
			return nil, false
		},
	})
}
