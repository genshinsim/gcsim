package spine

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("serpent spine", weapon)
	core.RegisterWeaponFunc("serpentspine", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	stacks := param["stacks"]
	c.Log.NewEvent(
		"serpent spine stack check",
		core.LogWeaponEvent,
		char.CharIndex(),
		"params", param,
	)
	buff := 0.05 + float64(r)*.01
	active := false

	c.Events.Subscribe(core.OnInitialize, func(args ...interface{}) bool {
		active = c.ActiveChar == char.CharIndex()
		return true
	}, fmt.Sprintf("spine-%v", char.Name()))

	c.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		if c.ActiveChar == char.CharIndex() {
			active = true
		} else {
			active = false
			//update stacks; duration is not reset yet by sim
			stacks = stacks + int(c.ActiveDuration/240)
			if stacks > 5 {
				stacks = 5
			}
		}
		return false
	}, fmt.Sprintf("spine-%v", char.Name()))

	c.Events.Subscribe(core.OnCharacterHurt, func(args ...interface{}) bool {
		if c.ActiveChar != char.CharIndex() {
			return false
		}
		stacks--
		if stacks < 0 {
			stacks = 0
		}
		return false
	}, fmt.Sprintf("spine-%v", char.Name()))

	val := make([]float64, core.EndStatType)
	char.AddMod(core.CharStatMod{
		Key:    "spine",
		Expiry: -1,
		Amount: func() ([]float64, bool) {
			//if active, then stacks = stacks + active dur
			//other wise it's just number of stacks
			x := stacks
			if active {
				x = stacks + int(c.ActiveDuration/240)
			}
			if x > 5 {
				x = 5
			}
			val[core.DmgP] = buff * float64(x)
			return val, true
		},
	})
	return "serpentspine"
}
