package hamayumi

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("hamayumi", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {

	nm := .12 + .04*float64(r)
	ca := .09 + .03*float64(r)

	char.AddPreDamageMod(coretype.PreDamageMod{
		Key:    "hamayumi",
		Expiry: -1,
		Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
			val := make([]float64, core.EndStatType)
			if atk.Info.AttackTag == coretype.AttackTagNormal {
				val[core.DmgP] = nm
				if char.CurrentEnergy() == char.MaxEnergy() {
					val[core.DmgP] = nm * 2
				}
				return val, true
			}

			if atk.Info.AttackTag == coretype.AttackTagExtra {
				val[core.DmgP] = ca
				if char.CurrentEnergy() == char.MaxEnergy() {
					val[core.DmgP] = ca * 2
				}
				return val, true
			}
			return nil, false
		},
	})

	return "hamayumi"
}
