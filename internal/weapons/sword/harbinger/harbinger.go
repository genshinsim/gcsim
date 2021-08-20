package harbinger

import (
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("harbinger of dawn", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {

	m := make([]float64, core.EndStatType)
	m[core.CR] = .105 + .035*float64(r)
	char.AddMod(core.CharStatMod{
		Key:    "harbinger",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, char.HP()/char.MaxHP() >= 0.9
		},
	})

}
