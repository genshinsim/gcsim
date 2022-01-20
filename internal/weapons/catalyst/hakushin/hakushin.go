package hakushin

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("hakushinring", weapon)
	core.RegisterWeaponFunc("hakushin", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {

	e := 10 + float64(r-1)*2.5
	e = e / 100
	m := make([]float64, core.EndStatType)
	m[core.ElectroP] = e
	hrfunc := func(ele core.EleType, key string) func(args ...interface{}) bool {
		return func(args ...interface{}) bool {
			if c.ActiveChar != char.CharIndex() {
				return false
			}
			m[ele] = e

			for _, char := range c.Chars {
				char.AddMod(core.CharStatMod{
					Key: "hakushin-passive-" + key,
					Amount: func() ([]float64, bool) {
						return m, true
					},
					Expiry: c.F + 6*60,
				})
			}
			return false
		}
	}

	c.Events.Subscribe(core.OnCrystallizeElectro, hrfunc(core.Geo, "hr-crystallize"), fmt.Sprintf("hakushin-ring-%v", char.Name()))
	c.Events.Subscribe(core.OnSwirlElectro, hrfunc(core.Anemo, "hr-swirl"), fmt.Sprintf("hakushin-ring-%v", char.Name()))
	c.Events.Subscribe(core.OnElectroCharged, hrfunc(core.Hydro, "hr-ec"), fmt.Sprintf("hakushin-ring-%v", char.Name()))
	c.Events.Subscribe(core.OnOverload, hrfunc(core.Pyro, "hr-ol"), fmt.Sprintf("hakushin-ring-%v", char.Name()))

	return "hakushinring"
}
