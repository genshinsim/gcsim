package deathmatch

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("deathmatch", weapon)
}

//If there are at least 2 opponents nearby, ATK is increased by 16% and DEF is increased by 16%.
//If there are fewer than 2 opponents nearby, ATK is increased by 24%.
func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	single := make([]float64, def.EndStatType)
	single[def.ATKP] = .12 + .04*float64(r)
	single[def.DEFP] = .12 + .04*float64(r)

	multiple := make([]float64, def.EndStatType)
	multiple[def.ATKP] = .18 + .06*float64(r)

	c.AddMod(def.CharStatMod{
		Key:    "deathmatch",
		Expiry: -1,
		Amount: func(a def.AttackTag) ([]float64, bool) {
			if len(s.Targets()) > 1 {
				return multiple, true
			}
			return single, true
		},
	})

}
