package oathsworneye

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("oathsworneye", weapon)
	core.RegisterWeaponFunc("oathsworn", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	val := make([]float64, core.EndStatType)
	val[core.ER] = 0.18 + 0.06*float64(r)
	c.Subscribe(core.PreSkill, func(args ...interface{}) bool {
		if c.ActiveChar != char.Index() {
			return false
		}
		char.AddMod(coretype.CharStatMod{
			Key: "oathsworn",
			Amount: func() ([]float64, bool) {
				return val, true
			},
			Expiry: c.Frame + 10*60,
		})
		return false
	}, fmt.Sprintf("oathsworn-%v", char.Name()))
	return "oathsworneye"
}
