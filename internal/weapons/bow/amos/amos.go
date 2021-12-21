package amos

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("amos' bow", weapon)
	core.RegisterWeaponFunc("amosbow", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {

	dmgpers := 0.06 + 0.02*float64(r)

	m := make([]float64, core.EndStatType)
	m[core.DmgP] = 0.09 + 0.03*float64(r)
	char.AddMod(core.CharStatMod{
		Key: "amos",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, a == core.AttackTagNormal || a == core.AttackTagExtra
		},
		Expiry: -1,
	})

	char.AddPreDamageMod(core.PreDamageMod{
		Key: "amos",
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			v := make([]float64, core.EndStatType)
			if atk.Info.AttackTag != core.AttackTagNormal && atk.Info.AttackTag != core.AttackTagExtra {
				return v, false
			}
			travel := float64(c.F-atk.SourceFrame) / 60
			stacks := int(travel / 0.1)
			if stacks > 5 {
				stacks = 5
			}
			v[core.DmgP] = dmgpers * float64(stacks)
			return v, true
		},
		Expiry: -1,
	})

}
