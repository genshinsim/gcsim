package festering

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("festering desire", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {

	m := make([]float64, core.EndStatType)
	m[core.CR] = .045 + .015*float64(r)
	m[core.DmgP] = .12 + 0.04*float64(r)
	c.AddMod(core.CharStatMod{
		Key:    "festering",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, a == core.AttackTagElementalArt
		},
	})

}
