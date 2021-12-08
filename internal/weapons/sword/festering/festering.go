package festering

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("festering desire", weapon)
	core.RegisterWeaponFunc("festeringdesire", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {

	var m [core.EndStatType]float64
	m[core.CR] = .045 + .015*float64(r)
	m[core.DmgP] = .12 + 0.04*float64(r)
	char.AddMod(core.CharStatMod{
		Key:    "festering",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([core.EndStatType]float64, bool) {
			return m, a == core.AttackTagElementalArt
		},
	})

}
