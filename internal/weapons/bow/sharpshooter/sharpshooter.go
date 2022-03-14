package sharpshooter

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("sharpshooter's oath", weapon)
	core.RegisterWeaponFunc("sharpshootersoath", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	dmg := 0.18 + float64(r)*0.06

	char.AddPreDamageMod(coretype.PreDamageMod{
		Key:    "sharpshooter",
		Expiry: -1,
		Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
			m := make([]float64, core.EndStatType)
			if atk.Info.HitWeakPoint {
				m[core.DmgP] = dmg
				return m, true
			}
			return nil, false
		},
	})

	return "sharpshootersoath"
}
