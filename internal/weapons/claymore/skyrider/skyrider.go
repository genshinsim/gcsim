package skyrider

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("skyrider greatsword", weapon)
	core.RegisterWeaponFunc("skyridergreatsword", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {

	atk := 0.05 + float64(r)*0.01
	stacks := 0
	icd := 0
	duration := 0

	c.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*coretype.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return false
		}
		if atk.Info.AttackTag != coretype.AttackTagNormal && atk.Info.AttackTag != coretype.AttackTagExtra {
			return false
		}
		if icd > c.Frame {
			return false
		}
		if duration < c.Frame {
			stacks = 0
		}

		stacks++
		if stacks > 4 {
			stacks = 4
		}
		icd = c.Frame + 30
		return false
	}, fmt.Sprintf("skyrider-greatsword-%v", char.Name()))

	val := make([]float64, core.EndStatType)
	char.AddMod(coretype.CharStatMod{
		Key:    "skyrider",
		Expiry: -1,
		Amount: func() ([]float64, bool) {
			if duration > c.Frame {
				val[core.ATKP] = atk * float64(stacks)
				return val, true
			}
			stacks = 0
			return nil, false
		},
	})
	return "skyridergreatsword"
}
