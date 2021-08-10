package harbinger

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("harbinger of dawn", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {

	m := make([]float64, core.EndStatType)
	m[core.CR] = .105 + .035*float64(r)
	c.AddMod(core.CharStatMod{
		Key:    "harbinger",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, c.HP()/c.MaxHP() >= 0.9
		},
	})

}
