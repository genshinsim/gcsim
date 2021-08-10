package grasscutter

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("grasscutter's light", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	atk := .21 + .07*float64(r)
	max := 0.7 + 0.1*float64(r)

	val := make([]float64, core.EndStatType)

	//ATK increased by 28% of Energy Recharge over the base 100%. You can gain a maximum bonus of 80% ATK.
	c.AddMod(core.CharStatMod{
		Key:          "grasscutter",
		Expiry:       -1,
		AffectedStat: core.ATKP, //this to prevent infinite loop when we ask to calculate ER
		Amount: func(a core.AttackTag) ([]float64, bool) {
			er := c.Stat(core.ER)
			bonus := atk * er
			if bonus > max {
				bonus = max
			}
			val[core.ATKP] = bonus
			return val, true
		},
	})

	erval := make([]float64, core.EndStatType)
	erval[core.ER] = .25 + .05*float64(r)

	//Gain 30% Energy Recharge for 12s after using an Elemental Burst.
	s.AddEventHook(func(s core.Sim) bool {
		if s.ActiveCharIndex() != c.CharIndex() {
			return false
		}
		c.AddMod(core.CharStatMod{
			Key:    "grasscutter-er",
			Expiry: 720,
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return erval, true
			},
		})

		return false
	}, fmt.Sprintf("grasscutter-%v", c.Name()), core.PostBurstHook)

}
