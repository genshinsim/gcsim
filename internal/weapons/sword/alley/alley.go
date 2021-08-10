package alley

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("the alley flash", weapon)
}

//Upon damaging an opponent, increases CRIT Rate by 8/10/12/14/16%. Max 5 stacks. A CRIT Hit removes all stacks.
func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {

	lockout := -1

	s.AddOnHurt(func(s core.Sim) {
		lockout = s.Frame() + 300
	})

	m := make([]float64, core.EndStatType)
	m[core.DmgP] = 0.09 + 0.03*float64(r)
	c.AddMod(core.CharStatMod{
		Key: "royal",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, lockout < s.Frame()
		},
		Expiry: -1,
	})

}
