package skyrider

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("skyrider sword", weapon)
	core.RegisterWeaponFunc("skyridersword", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {

	expiry := 0

	val := make([]float64, core.EndStatType)
	val[core.ATKP] = 0.09 + 0.03*float64(r)
	char.AddMod(coretype.CharStatMod{
		Key:    "skyrider",
		Expiry: -1,
		Amount: func() ([]float64, bool) {
			return val, expiry > c.Frame
		},
	})

	c.Subscribe(core.PostBurst, func(args ...interface{}) bool {
		if c.ActiveChar != char.Index() {
			return false
		}
		expiry = c.Frame + 900
		return false
	}, fmt.Sprintf("skyrider-sword-%v", char.Name()))

	return "skyridersword"
}
