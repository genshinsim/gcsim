package thrilling

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("thrilling tales of dragon slayers", weapon)
	core.RegisterWeaponFunc("thrillingtalesofdragonslayers", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	last := 0
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
			if last != 0 && c.F-last < 1200 {
				return false
			}
			//trigger buff if not on cd

			last = c.F
			expiry := c.F + 600

			active := c.Chars[c.ActiveChar]
			active.AddMod(core.CharStatMod{
				Key: "thrilling tales",
				Amount: func(a core.AttackTag) ([]float64, bool) {
					return m, expiry > c.F
				},
				Expiry: -1,
			})

			c.Log.Debugw("ttds activated", "frame", c.F, "event", core.LogWeaponEvent, "char", active.CharIndex(), "expiry", expiry)
		}

		return false
	}, fmt.Sprintf("thrilling-%v", char.Name()))

	return "thrillingtalesofdragonslayers"
}
