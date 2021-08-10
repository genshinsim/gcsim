package skyrider

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("skyrider sword", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {

	expiry := 0

	val := make([]float64, core.EndStatType)
	val[core.ATKP] = 0.09 + 0.03*float64(r)
	c.AddMod(core.CharStatMod{
		Key:    "skyrider",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return val, expiry > s.Frame()
		},
	})

	s.AddEventHook(func(s core.Sim) bool {
		if s.ActiveCharIndex() != c.CharIndex() {
			return false
		}
		expiry = s.Frame() + 900
		return false
	}, fmt.Sprintf("skyrider-sword-%v", c.Name()), core.PreBurstHook)

}
