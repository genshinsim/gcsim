package amenoma

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("amenoma kageuchi", weapon)
	core.RegisterWeaponFunc("amenomakageuchi", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {

	seeds := make([]int, 3) //keep track the seeds
	refund := 4.5 + 1.5*float64(r)
	icd := 0

	c.Subscribe(core.PostSkill, func(args ...interface{}) bool {
		if c.ActiveChar != char.Index() {
			return false
		}
		// add 1 seed
		if icd > c.Frame {
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

		seeds[index] = c.Frame + 30*60

		c.Log.NewEvent("amenoma proc'd", coretype.LogWeaponEvent, char.Index(), "index", index, "seeds", seeds)

		icd = c.Frame + 300 //5 seconds

		return false
	}, fmt.Sprintf("amenoma-skill-%v", char.Name()))

	c.Subscribe(core.PostBurst, func(args ...interface{}) bool {
		if c.ActiveChar != char.Index() {
			return false
		}
		count := 0
		for i, v := range seeds {
			if v > c.Frame {
				count++
			}
			seeds[i] = 0
		}
		if count == 0 {
			return false
		}
		//regen energy after 2 seconds
		char.AddTask(func() {
			char.AddEnergy("amenoma", refund*float64(count))
		}, "amenoma-regen", 120+60) //added 1 extra sec for burst animation but who knows if this is true

		return false
	}, fmt.Sprintf("amenoma-burst-%v", char.Name()))
	return "amenomakageuchi"
}
