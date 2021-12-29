package vortex

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("vortex vanquisher", weapon)
	core.RegisterWeaponFunc("vortexvanquisher", weapon)
}

//Increases DMG against enemies affected by Hydro or Electro by 20/24/28/32/36%.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {

	shd := .15 + float64(r)*.05
	c.Shields.AddBonus(func() float64 {
		return shd
	})

	stacks := 0
	icd := 0
	duration := 0

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if icd > c.F {
			return false
		}
		if duration < c.F {
			stacks = 0
		}
		stacks++
		if stacks > 5 {
			stacks = 0
		}
		icd = c.F + 18
		return false

	}, fmt.Sprintf("vortex-%v", char.Name()))

	atk := 0.03 + 0.01*float64(r)

	char.AddMod(core.CharStatMod{
		Key:    "vortex",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			val := make([]float64, core.EndStatType)
			if duration > c.F {
				val[core.ATKP] = atk * float64(stacks)
				if c.Shields.IsShielded() {
					val[core.ATKP] *= 2
				}
				return val, true
			}
			stacks = 0
			return nil, false
		},
	})
	return "vortexvanquisher"
}
