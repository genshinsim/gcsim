package hamayumi

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("hamayumi", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {

	nm := .12 + .04*float64(r)
	ca := .09 + .03*float64(r)

	val := make([]float64, core.EndStatType)
	c.AddMod(core.CharStatMod{
		Key:    "hamayumi",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if a == core.AttackTagNormal {
				val[core.DmgP] = nm
				if c.CurrentEnergy() == c.MaxEnergy() {
					val[core.DmgP] = nm * 2
				}
				return val, true
			}

			if a == core.AttackTagExtra {
				val[core.DmgP] = ca
				if c.CurrentEnergy() == c.MaxEnergy() {
					val[core.DmgP] = ca * 2
				}
				return val, true
			}
			return nil, false
		},
	})
}
