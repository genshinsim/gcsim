package sharpshooter

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("sharpshooter's oath", weapon)
	core.RegisterWeaponFunc("sharpshootersoath", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	dmg := 0.18 + float64(r)*0.06

	char.AddPreDamageMod(core.PreDamageMod{
		Key:    "sharpshooter",
		Expiry: -1,
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			m := make([]float64, core.EndStatType)
			if atk.Info.HitWeakPoint {
				m[core.DmgP] = dmg
				return m, true
			}
			return nil, false
		},
	})

}
