package harbinger

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("harbinger of dawn", weapon)
	core.RegisterWeaponFunc("harbingerofdawn", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {

	m := make([]float64, core.EndStatType)
	m[core.CR] = .105 + .035*float64(r)
	char.AddMod(core.CharStatMod{
		Key:    "harbinger",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, char.HP()/char.MaxHP() >= 0.9
		},
	})
	return "harbingerofdawn"
}
