package royal

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("royal spear", weapon)
}

//Upon damaging an opponent, increases CRIT Rate by 8/10/12/14/16%. Max 5 stacks. A CRIT Hit removes all stacks.
func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	stacks := 0

	s.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		if crit {
			stacks = 0
		} else {
			stacks++
			if stacks > 5 {
				stacks = 5
			}
		}
	}, fmt.Sprintf("royal-spear-%v", c.Name()))

	rate := 0.06 + float64(r)*0.02
	m := make([]float64, def.EndStatType)
	c.AddMod(def.CharStatMod{
		Key: "royal",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			m[def.CR] = float64(stacks) * rate
			return m, true
		},
		Expiry: -1,
	})

}
