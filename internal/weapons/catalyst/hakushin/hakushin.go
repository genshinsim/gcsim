package hakushin

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("hakushin ring", weapon)
	core.RegisterWeaponFunc("hakushin", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {

	expiry := 0
	e := 10 + float64(r)*2.5
	m := make([]float64, core.EndStatType)
	m[core.ElectroP] = e
	hrfunc := func(ele core.EleType, key string) func(args ...interface{}) bool {
		return func(args ...interface{}) bool {
			if c.ActiveChar != char.CharIndex() {
				return false
			}
			m[ele] = e
			expiry = c.F + 6*60
			return false
		}
	}

	for _, char := range c.Chars {
		char.AddMod(core.CharStatMod{
			Key: "hakushin-passive",
			Amount: func() ([]float64, bool) {
				return m, expiry > c.F
			},
			Expiry: -1,
		})
	}

	c.Events.Subscribe(core.OnCrystallizeElectro, hrfunc(core.Geo, "hr-crystallize"), fmt.Sprintf("hakushin-ring-%v", char.Name()))
	c.Events.Subscribe(core.OnSwirlElectro, hrfunc(core.Anemo, "hr-swirl"), fmt.Sprintf("hakushin-ring-%v", char.Name()))
	c.Events.Subscribe(core.OnElectroCharged, hrfunc(core.Hydro, "hr-ec"), fmt.Sprintf("hakushin-ring-%v", char.Name()))
	c.Events.Subscribe(core.OnOverload, hrfunc(core.Pyro, "hr-ol"), fmt.Sprintf("hakushin-ring-%v", char.Name()))

	return "hakushinring"
}
