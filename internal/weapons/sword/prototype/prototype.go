package prototype

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("prototype rancour", weapon)
	core.RegisterWeaponFunc("prototyperancour", weapon)
}

//After using an Elemental Skill, increases Normal and Charged Attack DMG by 8% for 12s. Max 2 stacks.
func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {

	expiry := 0
	per := 0.03 + 0.01*float64(r)
	stacks := 0
	icd := 0

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
		icd = c.Frame + 18
		if expiry < c.Frame {
			stacks = 0
		}
		stacks++
		if stacks > 4 {
			stacks = 4
		}
		expiry = c.Frame + 360
		return false
	}, fmt.Sprintf("prototype-rancour-%v", char.Name()))

	char.AddMod(coretype.CharStatMod{
		Key:    "prototype",
		Expiry: -1,
		Amount: func() ([]float64, bool) {
			val := make([]float64, core.EndStatType)
			if expiry < c.Frame {
				stacks = 0
				return nil, false
			}
			val[core.ATKP] = per * float64(stacks)
			val[core.DEFP] = per * float64(stacks)
			return val, true
		},
	})

	return "prototyperancour"

}
