package skyrider

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("skyrider greatsword", weapon)
	core.RegisterWeaponFunc("skyridergreatsword", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {

	atk := 0.05 + float64(r)*0.01
	stacks := 0
	icd := 0
	duration := 0

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if atk.Info.AttackTag != core.AttackTagNormal && atk.Info.AttackTag != core.AttackTagExtra {
			return false
		}
		if icd > c.F {
			return false
		}
		if duration < c.F {
			stacks = 0
		}

		stacks++
		if stacks > 4 {
			stacks = 4
		}
		icd = c.F + 30
		return false
	}, fmt.Sprintf("skyrider-greatsword-%v", char.Name()))

	val := make([]float64, core.EndStatType)
	char.AddMod(core.CharStatMod{
		Key:    "skyrider",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if duration > c.F {
				val[core.ATKP] = atk * float64(stacks)
				return val, true
			}
			stacks = 0
			return nil, false
		},
	})

}
