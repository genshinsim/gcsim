package spine

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("serpent spine", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	stacks := param["stacks"]
	buff := 0.05 + float64(r)*.01
	active := false

	s.AddInitHook(func() {
		active = s.ActiveCharIndex() == c.CharIndex()
	})

	s.AddEventHook(func(s core.Sim) bool {
		if s.ActiveCharIndex() == c.CharIndex() {
			active = true
		} else {
			active = false
			//update stacks; duration is not reset yet by sim
			stacks = stacks + int(s.ActiveDuration()/240)
			if stacks > 5 {
				stacks = 5
			}
		}
		return false
	}, fmt.Sprintf("spine-%v", c.Name()), core.PostSwapHook)

	s.AddOnHurt(func(s core.Sim) {
		stacks--
		if stacks < 0 {
			stacks = 0
		}
	})

	val := make([]float64, core.EndStatType)
	c.AddMod(core.CharStatMod{
		Key:    "spine",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			//if active, then stacks = stacks + active dur
			//other wise it's just number of stacks
			x := stacks
			if active {
				x = stacks + int(s.ActiveDuration()/240)
			}
			if x > 5 {
				x = 5
			}
			val[core.DmgP] = buff * float64(x)
			return val, true
		},
	})
}
