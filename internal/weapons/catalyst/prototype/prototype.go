package prototype

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("prototype amber", weapon)
}

//Using an Elemental Burst regenerates 4/4.5/5/5.5/6 Energy every 2s for 6s. All party members
//will regenerate 4/4.5/5/5.5/6% HP every 2s for this duration.
func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {

	e := 3.5 + float64(r)*0.5

	s.AddEventHook(func(s core.Sim) bool {

		for i := 120; i <= 360; i += 120 {
			c.AddTask(func() {
				for _, char := range s.Characters() {
					char.AddEnergy(e)
				}
				s.HealAllPercent(e / 100.0)
			}, "recharge", i)
		}

		return false
	}, fmt.Sprintf("prototype-amber-%v", c.Name()), core.PostBurstHook)
}
