package skyrider

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("skyrider sword", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {

	expiry := 0

	val := make([]float64, def.EndStatType)
	val[def.ATKP] = 0.09 + 0.03*float64(r)
	c.AddMod(def.CharStatMod{
		Key:    "skyrider",
		Expiry: -1,
		Amount: func(a def.AttackTag) ([]float64, bool) {
			return val, expiry > s.Frame()
		},
	})

	s.AddEventHook(func(s def.Sim) bool {
		if s.ActiveCharIndex() != c.CharIndex() {
			return false
		}
		expiry = s.Frame() + 900
		return false
	}, fmt.Sprintf("skyrider-sword-%v", c.Name()), def.PreBurstHook)

}
