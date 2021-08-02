package grasscutter

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("grasscutter's light", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	atk := .21 + .07*float64(r)
	max := 0.7 + 0.1*float64(r)

	val := make([]float64, def.EndStatType)

	//ATK increased by 28% of Energy Recharge over the base 100%. You can gain a maximum bonus of 80% ATK.
	c.AddMod(def.CharStatMod{
		Key:          "grasscutter",
		Expiry:       -1,
		AffectedStat: def.ATKP, //this to prevent infinite loop when we ask to calculate ER
		Amount: func(a def.AttackTag) ([]float64, bool) {
			er := c.Stat(def.ER)
			bonus := atk * er
			if bonus > max {
				bonus = max
			}
			val[def.ATKP] = bonus
			return val, true
		},
	})

	erval := make([]float64, def.EndStatType)
	erval[def.ER] = .25 + .05*float64(r)

	//Gain 30% Energy Recharge for 12s after using an Elemental Burst.
	s.AddEventHook(func(s def.Sim) bool {
		if s.ActiveCharIndex() != c.CharIndex() {
			return false
		}
		c.AddMod(def.CharStatMod{
			Key:    "grasscutter-er",
			Expiry: 720,
			Amount: func(a def.AttackTag) ([]float64, bool) {
				return erval, true
			},
		})

		return false
	}, fmt.Sprintf("grasscutter-%v", c.Name()), def.PostBurstHook)

}
