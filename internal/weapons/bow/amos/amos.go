package amos

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("amos' bow", weapon)
	core.RegisterWeaponFunc("amosbow", weapon)
	core.RegisterWeaponFunc("amos", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {

	dmgpers := 0.06 + 0.02*float64(r)

	m := make([]float64, core.EndStatType)
	// m[core.DmgP] = 0.09 + 0.03*float64(r)
	flat := 0.09 + 0.03*float64(r)

	char.AddPreDamageMod(coretype.PreDamageMod{
		Key: "amos",
		Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
			if atk.Info.AttackTag != coretype.AttackTagNormal && atk.Info.AttackTag != coretype.AttackTagExtra {
				return nil, false
			}
			m[core.DmgP] = flat
			travel := float64(c.Frame-atk.Snapshot.SourceFrame) / 60
			stacks := int(travel / 0.1)
			if stacks > 5 {
				stacks = 5
			}
			m[core.DmgP] += dmgpers * float64(stacks)
			return m, true
		},
		Expiry: -1,
	})

	return "amosbow"
}
