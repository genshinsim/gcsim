package royal

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("royal bow", weapon)
	core.RegisterWeaponFunc("royal grimore", weapon)
	core.RegisterWeaponFunc("royal greatsword", weapon)
	core.RegisterWeaponFunc("royal spear", weapon)
	core.RegisterWeaponFunc("royal longsword", weapon)
}

//Upon damaging an opponent, increases CRIT Rate by 8/10/12/14/16%. Max 5 stacks. A CRIT Hit removes all stacks.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	stacks := 0

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		crit := args[3].(bool)
		if crit {
			stacks = 0
		} else {
			stacks++
			if stacks > 5 {
				stacks = 5
			}
		}
		return false
	}, fmt.Sprintf("royal-%v", char.Name()))

	rate := 0.06 + float64(r)*0.02
	m := make([]float64, core.EndStatType)
	char.AddMod(core.CharStatMod{
		Key: "royal",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			m[core.CR] = float64(stacks) * rate
			return m, true
		},
		Expiry: -1,
	})

}
