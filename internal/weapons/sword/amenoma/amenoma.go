package amenoma

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("amenoma kageuchi", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {

	seeds := make([]int, 3) //keep track the seeds

	icd := 0

	s.AddEventHook(func(s def.Sim) bool {
		// add 1 seed
		if icd > s.Frame() {
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

		seeds[index] = s.Frame() + 30*60
		icd = s.Frame() + 300 //5 seconds

		return false
	}, fmt.Sprintf("amenoma-skill-%v", c.Name()), def.PostSkillHook)

	s.AddEventHook(func(s def.Sim) bool {
		if s.ActiveCharIndex() != c.CharIndex() {
			return false
		}
		count := 0
		for i, v := range seeds {
			if v > s.Frame() {
				count++
			}
			seeds[i] = 0
		}
		if count == 0 {
			return false
		}
		//regen energy after 2 seconds
		c.AddTask(func() {
			c.AddEnergy(6 * float64(count))
		}, "amenoma-regen", 120+60) //added 1 extra sec for burst animation but who knows if this is true

		return false
	}, fmt.Sprintf("amenoma-burst-%v", c.Name()), def.PostBurstHook)

}
