package prototype

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("prototype rancour", weapon)
	core.RegisterWeaponFunc("prototyperancour", weapon)
}

//After using an Elemental Skill, increases Normal and Charged Attack DMG by 8% for 12s. Max 2 stacks.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) {

	expiry := 0
	per := 0.03 + 0.01*float64(r)
	stacks := 0
	icd := 0

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
		icd = c.F + 18
		if expiry < c.F {
			stacks = 0
		}
		stacks++
		if stacks > 4 {
			stacks = 4
		}
		expiry = c.F + 360
		return false
	}, fmt.Sprintf("prototype-rancour-%v", char.Name()))

	char.AddMod(core.CharStatMod{
		Key:    "prototype",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			val := make([]float64, core.EndStatType)
			if expiry < c.F {
				stacks = 0
				return nil, false
			}
			val[core.ATKP] = per * float64(stacks)
			val[core.DEFP] = per * float64(stacks)
			return val, true
		},
	})

}
