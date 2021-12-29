package skyrider

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("skyrider sword", weapon)
	core.RegisterWeaponFunc("skyridersword", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {

	expiry := 0

	val := make([]float64, core.EndStatType)
	val[core.ATKP] = 0.09 + 0.03*float64(r)
	char.AddMod(core.CharStatMod{
		Key:    "skyrider",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return val, expiry > c.F
		},
	})

	c.Events.Subscribe(core.PostBurst, func(args ...interface{}) bool {
		if c.ActiveChar != char.CharIndex() {
			return false
		}
		expiry = c.F + 900
		return false
	}, fmt.Sprintf("skyrider-sword-%v", char.Name()))

	return "skyridersword"
}
