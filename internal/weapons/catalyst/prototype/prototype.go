package prototype

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("prototype amber", weapon)
	core.RegisterWeaponFunc("prototypeamber", weapon)
}

//Using an Elemental Burst regenerates 4/4.5/5/5.5/6 Energy every 2s for 6s. All party members
//will regenerate 4/4.5/5/5.5/6% HP every 2s for this duration.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {

	e := 3.5 + float64(r)*0.5

	c.Events.Subscribe(core.PostBurst, func(args ...interface{}) bool {

		for i := 120; i <= 360; i += 120 {
			char.AddTask(func() {
				char.AddEnergy(e)
				c.Health.HealAllPercent(char.CharIndex(), e/100.0)
			}, "recharge", i)
		}

		return false
	}, fmt.Sprintf("prototype-amber-%v", char.Name()))

	return "prototypeamber"
}
