package lithic

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("lithic blade", weapon)
}

//For every character in the party who hails from Liyue, the character who equips this
//weapon gains 6/7/8/9//10% ATK increase and 2/3/4/5/6% CRIT Rate increase.
func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {

	stacks := 0

	s.AddInitHook(func() {
		for _, char := range s.Characters() {
			if char.Zone() == def.ZoneLiyue {
				stacks++
			}
		}
	})

	val := make([]float64, def.EndStatType)
	val[def.CR] = (0.02 + float64(r)*0.01) * float64(stacks)
	val[def.ATKP] = (0.06 + float64(r)*0.01) * float64(stacks)

	c.AddMod(def.CharStatMod{
		Key:    "lithic",
		Expiry: -1,
		Amount: func(a def.AttackTag) ([]float64, bool) {
			return val, true
		},
	})
}
