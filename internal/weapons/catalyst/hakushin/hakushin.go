package hakushin

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("hakushinring", weapon)
	core.RegisterWeaponFunc("hakushin ring", weapon)
	core.RegisterWeaponFunc("hakushin", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {

	m := make([]float64, core.EndStatType)
	dmg := .075 + float64(r)*.025

	hrfunc := func(ele coretype.EleType, key string) func(args ...interface{}) bool {
		icd := -1
		return func(args ...interface{}) bool {
			ae := args[1].(*coretype.AttackEvent)

			if c.ActiveChar != char.Index() {
				return false
			}
			if ae.Info.ActorIndex != char.Index() {
				return false
			}

			// do not overwrite mod if same frame
			if c.Frame < icd {
				return false
			}
			icd = c.Frame + 1

			for _, char := range c.Chars {
				if char.Ele() != core.Electro && char.Ele() != ele {
					continue
				}
				this := char
				char.AddMod(coretype.CharStatMod{
					Key: "hakushin-passive",
					Amount: func() ([]float64, bool) {

						m[core.PyroP] = 0
						m[core.HydroP] = 0
						m[coretype.CryoP] = 0
						m[core.ElectroP] = 0
						m[core.AnemoP] = 0
						m[core.GeoP] = 0
						m[core.DendroP] = 0
						m[coretype.EleToDmgP(this.Ele())] = dmg

						return m, true
					},
					Expiry: c.Frame + 6*60,
				})
			}

			c.Log.NewEvent("hakushin proc'd", coretype.LogWeaponEvent, char.Index(), "trigger", key, "expiring", c.Frame+6*60)

			return false
		}
	}

	c.Subscribe(core.OnCrystallizeElectro, hrfunc(core.Geo, "hr-crystallize"), fmt.Sprintf("hakushin-ring-%v", char.Name()))
	c.Subscribe(coretype.OnSwirlElectro, hrfunc(core.Anemo, "hr-swirl"), fmt.Sprintf("hakushin-ring-%v", char.Name()))
	c.Subscribe(core.OnElectroCharged, hrfunc(core.Hydro, "hr-ec"), fmt.Sprintf("hakushin-ring-%v", char.Name()))
	c.Subscribe(core.OnOverload, hrfunc(core.Pyro, "hr-ol"), fmt.Sprintf("hakushin-ring-%v", char.Name()))
	c.Subscribe(core.OnSuperconduct, hrfunc(coretype.Cryo, "hr-sc"), fmt.Sprintf("hakushin-ring-%v", char.Name()))

	return "hakushinring"
}
