package harbinger

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("harbinger of dawn", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {

	m := make([]float64, def.EndStatType)
	m[def.CR] = .105 + .035*float64(r)
	c.AddMod(def.CharStatMod{
		Key:    "harbinger",
		Expiry: -1,
		Amount: func(a def.AttackTag) ([]float64, bool) {
			return m, c.HP()/c.MaxHP() >= 0.9
		},
	})

}
