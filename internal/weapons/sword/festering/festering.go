package festering

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("festering desire", weapon)
	core.RegisterWeaponFunc("festeringdesire", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {

	m := make([]float64, core.EndStatType)
	m[core.CR] = .045 + .015*float64(r)
	m[core.DmgP] = .12 + 0.04*float64(r)
	char.AddPreDamageMod(coretype.PreDamageMod{
		Key:    "festering",
		Expiry: -1,
		Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
			switch atk.Info.AttackTag {
			case core.AttackTagElementalArt, core.AttackTagElementalArtHold:
				return m, true
			}
			return nil, false
		},
	})
	return "festeringdesire"
}
