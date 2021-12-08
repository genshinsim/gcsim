package deathmatch

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("deathmatch", weapon)
}

//If there are at least 2 opponents nearby, ATK is increased by 16% and DEF is increased by 16%.
//If there are fewer than 2 opponents nearby, ATK is increased by 24%.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) {

	var multiple [core.EndStatType]float64
	multiple[core.ATKP] = .12 + .04*float64(r)
	multiple[core.DEFP] = .12 + .04*float64(r)

	var single [core.EndStatType]float64
	single[core.ATKP] = .18 + .06*float64(r)

	char.AddMod(core.CharStatMod{
		Key:    "deathmatch",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([core.EndStatType]float64, bool) {
			if len(c.Targets) > 1 {
				return multiple, true
			}
			return single, true
		},
	})

}
