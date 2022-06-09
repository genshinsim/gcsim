package thrilling

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("thrilling tales of dragon slayers", weapon)
	core.RegisterWeaponFunc("thrillingtalesofdragonslayers", weapon)
	core.RegisterWeaponFunc("ttds", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	cd := -1
	isActive := false

	c.Events.Subscribe(core.OnInitialize, func(args ...interface{}) bool {
		isActive = c.ActiveChar == char.CharIndex()
		return true
	}, fmt.Sprintf("thrilling-%v", char.Name()))

	m := make([]float64, core.EndStatType)
	m[core.ATKP] = .18 + float64(r)*0.06

	c.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		if !isActive && c.ActiveChar == char.CharIndex() {
			//swapped to current char
			isActive = true
			return false
		}

		//swap from current char to new char
		if isActive && c.ActiveChar != char.CharIndex() {
			isActive = false

			//do nothing if off cd
			if c.F < cd {
				return false
			}
			//trigger buff if not on cd
			cd = c.F + 60*20
			expiry := c.F + 60*10

			active := c.Chars[c.ActiveChar]
			active.AddMod(core.CharStatMod{
				Key: "thrilling tales",
				Amount: func() ([]float64, bool) {
					return m, true
				},
				Expiry: expiry,
			})

			c.Log.NewEvent("ttds activated", core.LogWeaponEvent, active.CharIndex(), "expiry", expiry)
		}

		return false
	}, fmt.Sprintf("thrilling-%v", char.Name()))

	return "thrillingtalesofdragonslayers"
}
