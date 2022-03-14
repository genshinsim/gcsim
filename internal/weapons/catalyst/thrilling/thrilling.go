package thrilling

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("thrilling tales of dragon slayers", weapon)
	core.RegisterWeaponFunc("thrillingtalesofdragonslayers", weapon)
	core.RegisterWeaponFunc("ttds", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	cd := -1
	isActive := false

	c.Subscribe(core.OnInitialize, func(args ...interface{}) bool {
		isActive = c.ActiveChar == char.Index()
		return true
	}, fmt.Sprintf("thrilling-%v", char.Name()))

	m := make([]float64, core.EndStatType)
	m[core.ATKP] = .18 + float64(r)*0.06

	c.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		if !isActive && c.ActiveChar == char.Index() {
			//swapped to current char
			isActive = true
			return false
		}

		//swap from current char to new char
		if isActive && c.ActiveChar != char.Index() {
			isActive = false

			//do nothing if off cd
			if c.Frame < cd {
				return false
			}
			//trigger buff if not on cd
			cd = c.Frame + 60*20
			expiry := c.Frame + 60*10

			active := c.Chars[c.ActiveChar]
			active.AddMod(coretype.CharStatMod{
				Key: "thrilling tales",
				Amount: func() ([]float64, bool) {
					return m, true
				},
				Expiry: expiry,
			})

			c.Log.NewEvent("ttds activated", coretype.LogWeaponEvent, active.Index(), "expiry", expiry)
		}

		return false
	}, fmt.Sprintf("thrilling-%v", char.Name()))

	return "thrillingtalesofdragonslayers"
}
