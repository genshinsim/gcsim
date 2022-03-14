package grasscutter

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("grasscutter's light", weapon)
	core.RegisterWeaponFunc("engulfinglightning", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	atk := .21 + .07*float64(r)
	max := 0.7 + 0.1*float64(r)

	val := make([]float64, core.EndStatType)

	//ATK increased by 28% of Energy Recharge over the base 100%. You can gain a maximum bonus of 80% ATK.
	char.AddMod(coretype.CharStatMod{
		Key:          "grasscutter",
		Expiry:       -1,
		AffectedStat: core.ATKP, //this to prevent infinite loop when we ask to calculate ER
		Amount: func() ([]float64, bool) {
			er := char.Stat(core.ER)
			c.Log.NewEvent("cutter snapshot", coretype.LogWeaponEvent, char.Index(), "er", er)
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
	c.Subscribe(core.PreBurst, func(args ...interface{}) bool {
		if c.ActiveChar != char.Index() {
			return false
		}
		char.AddMod(coretype.CharStatMod{
			Key:    "grasscutter-er",
			Expiry: c.Frame + 720,
			Amount: func() ([]float64, bool) {
				return erval, true
			},
		})

		return false
	}, fmt.Sprintf("grasscutter-%v", char.Name()))
	return "engulfinglightning"
}
