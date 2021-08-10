package homa

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("staff of homa", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	//add on hit effect to sim?
	m := make([]float64, core.EndStatType)
	m[core.HPP] = 0.15 + float64(r)*0.05
	atkp := 0.006 + float64(r)*0.002
	lowhp := 0.008 + float64(r)*0.002

	c.AddMod(core.CharStatMod{
		Key: "homa hp bonus",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			per := atkp
			if c.HP()/c.MaxHP() <= 0.5 {
				per += lowhp
			}
			m[core.ATKP] = per * c.MaxHP()
			return m, true
		},
		Expiry: -1,
	})

}
