package hamayumi

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("hamayumi", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {

	nm := .12 + .04*float64(r)
	ca := .09 + .03*float64(r)

	val := make([]float64, def.EndStatType)
	c.AddMod(def.CharStatMod{
		Key:    "hamayumi",
		Expiry: -1,
		Amount: func(a def.AttackTag) ([]float64, bool) {
			if a == def.AttackTagNormal {
				val[def.DmgP] = nm
				if c.CurrentEnergy() == c.MaxEnergy() {
					val[def.DmgP] = nm * 2
				}
				return val, true
			}

			if a == def.AttackTagExtra {
				val[def.DmgP] = ca
				if c.CurrentEnergy() == c.MaxEnergy() {
					val[def.DmgP] = ca * 2
				}
				return val, true
			}
			return nil, false
		},
	})
}
