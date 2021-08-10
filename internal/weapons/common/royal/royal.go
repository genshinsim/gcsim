package royal

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("royal bow", weapon)
	combat.RegisterWeaponFunc("royal grimore", weapon)
	combat.RegisterWeaponFunc("royal greatsword", weapon)
	combat.RegisterWeaponFunc("royal spear", weapon)
	combat.RegisterWeaponFunc("royal longsword", weapon)
}

//Upon damaging an opponent, increases CRIT Rate by 8/10/12/14/16%. Max 5 stacks. A CRIT Hit removes all stacks.
func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	stacks := 0

	s.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
		if crit {
			stacks = 0
		} else {
			stacks++
			if stacks > 5 {
				stacks = 5
			}
		}
	}, fmt.Sprintf("royal-%v", c.Name()))

	rate := 0.06 + float64(r)*0.02
	m := make([]float64, core.EndStatType)
	c.AddMod(core.CharStatMod{
		Key: "royal",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			m[core.CR] = float64(stacks) * rate
			return m, true
		},
		Expiry: -1,
	})

}
