package alley

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("the alley flash", weapon)
}

//Upon damaging an opponent, increases CRIT Rate by 8/10/12/14/16%. Max 5 stacks. A CRIT Hit removes all stacks.
func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {

	lockout := -1

	s.AddOnHurt(func(s def.Sim) {
		lockout = s.Frame() + 300
	})

	m := make([]float64, def.EndStatType)
	m[def.DmgP] = 0.09 + 0.03*float64(r)
	c.AddMod(def.CharStatMod{
		Key: "royal",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			return m, lockout < s.Frame()
		},
		Expiry: -1,
	})

}
