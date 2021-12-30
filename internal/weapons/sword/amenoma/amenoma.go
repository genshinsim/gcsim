package amenoma

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("amenoma kageuchi", weapon)
	core.RegisterWeaponFunc("amenomakageuchi", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {

	seeds := make([]int, 3) //keep track the seeds
	refund := 4.5 + 1.5*float64(r)
	icd := 0

	c.Events.Subscribe(core.PostSkill, func(args ...interface{}) bool {
		if c.ActiveChar != char.CharIndex() {
			return false
		}
		// add 1 seed
		if icd > c.F {
			return false
		}
		// find oldest seed to overwrite
		index := 0
		old := seeds[0]

		for i, v := range seeds {
			if v < old {
				old = v
				index = i
			}
		}

		seeds[index] = c.F + 30*60

		c.Log.Debugw("amenoma proc'd", "frame", c.F, "event", core.LogWeaponEvent, "char", char.CharIndex(), "index", index, "seeds", seeds)

		icd = c.F + 300 //5 seconds

		return false
	}, fmt.Sprintf("amenoma-skill-%v", char.Name()))

	c.Events.Subscribe(core.PostBurst, func(args ...interface{}) bool {
		if c.ActiveChar != char.CharIndex() {
			return false
		}
		count := 0
		for i, v := range seeds {
			if v > c.F {
				count++
			}
			seeds[i] = 0
		}
		if count == 0 {
			return false
		}
		//regen energy after 2 seconds
		char.AddTask(func() {
			char.AddEnergy(refund * float64(count))
		}, "amenoma-regen", 120+60) //added 1 extra sec for burst animation but who knows if this is true

		return false
	}, fmt.Sprintf("amenoma-burst-%v", char.Name()))
	return "amenomakageuchi"
}
