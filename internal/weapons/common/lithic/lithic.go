package lithic

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("lithic spear", weapon)
	combat.RegisterWeaponFunc("lithic blade", weapon)
}

//For every character in the party who hails from Liyue, the character who equips this
//weapon gains 6/7/8/9//10% ATK increase and 2/3/4/5/6% CRIT Rate increase.
func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {

	stacks := 0

	s.AddInitHook(func() {
		for _, char := range s.Characters() {
			if char.Zone() == core.ZoneLiyue {
				stacks++
			}
		}
	})

	val := make([]float64, core.EndStatType)
	val[core.CR] = (0.02 + float64(r)*0.01) * float64(stacks)
	val[core.ATKP] = (0.06 + float64(r)*0.01) * float64(stacks)

	c.AddMod(core.CharStatMod{
		Key:    "lithic",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return val, true
		},
	})
}
