package primordial

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("primordial jade cutter", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	//add on hit effect to sim?
	m := make([]float64, core.EndStatType)
	m[core.HPP] = 0.15 + float64(r)*0.05
	atkp := 0.009 + float64(r)*0.003

	c.AddMod(core.CharStatMod{
		Key: "homa hp bonus",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			m[core.ATKP] = atkp * c.MaxHP()
			return m, true
		},
		Expiry: -1,
	})

}
