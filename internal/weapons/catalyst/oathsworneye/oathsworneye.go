package oathsworneye

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("oathsworneye", weapon)
	core.RegisterWeaponFunc("oathsworn", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	val := make([]float64, core.EndStatType)
	val[core.ER] = 0.18 + 0.06*float64(r)
	c.Events.Subscribe(core.PreSkill, func(args ...interface{}) bool {
		if c.ActiveChar != char.CharIndex() {
			return false
		}
		char.AddMod(core.CharStatMod{
			Key: "oathsworn",
			Amount: func() ([]float64, bool) {
				return val, true
			},
			Expiry: c.F + 10*60,
		})
		return false
	}, fmt.Sprintf("oathsworn-%v", char.Name()))
	return "oathsworneye"
}
