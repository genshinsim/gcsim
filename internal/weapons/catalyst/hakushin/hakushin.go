package hakushin

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("hakushinring", weapon)
	core.RegisterWeaponFunc("hakushin ring", weapon)
	core.RegisterWeaponFunc("hakushin", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {

	dmg := .075 + float64(r)*.025

	hrfunc := func(ele core.EleType, key string) func(args ...interface{}) bool {
		return func(args ...interface{}) bool {
			ae := args[1].(*core.AttackEvent)

			if c.ActiveChar != char.CharIndex() {
				return false
			}
			if ae.Info.ActorIndex != char.CharIndex() {
				return false
			}

			for _, char := range c.Chars {
				m := make([]float64, core.EndStatType)

				switch charEle := char.Ele(); charEle {
				case core.Electro, ele:
					m[core.EleToDmgP(charEle)] = dmg
				default:
					continue
				}

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
	c.Events.Subscribe(core.OnSuperconduct, hrfunc(core.Cryo, "hr-sc"), fmt.Sprintf("hakushin-ring-%v", char.Name()))

	return "hakushinring"
}
