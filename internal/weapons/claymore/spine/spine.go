package spine

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("serpent spine", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	stacks := param["stacks"]
	buff := 0.05 + float64(r)*.01
	active := false

	s.AddInitHook(func() {
		active = s.ActiveCharIndex() == c.CharIndex()
	})

	s.AddEventHook(func(s def.Sim) bool {
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
	}, fmt.Sprintf("spine-%v", c.Name()), def.PostSwapHook)

	s.AddOnHurt(func(s def.Sim) {
		stacks--
		if stacks < 0 {
			stacks = 0
		}
	})

	val := make([]float64, def.EndStatType)
	c.AddMod(def.CharStatMod{
		Key:    "spine",
		Expiry: -1,
		Amount: func(a def.AttackTag) ([]float64, bool) {
			//if active, then stacks = stacks + active dur
			//other wise it's just number of stacks
			x := stacks
			if active {
				x = stacks + int(s.ActiveDuration()/240)
			}
			if x > 5 {
				x = 5
			}
			val[def.DmgP] = buff * float64(x)
			return val, true
		},
	})
}
