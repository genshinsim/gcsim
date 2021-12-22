package memory

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("memory of dust", weapon)
	core.RegisterWeaponFunc("memoryofdust", weapon)
}

//Increases Shield Strength by 20/25/30/35/40%. Scoring hits on opponents increases ATK by 4/5/6/7/8% for 8s. Max 5 stacks.
//Can only occur once every 0.3s. While protected by a shield, this ATK increase effect is increased by 100%.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) {

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
		return false
	}, fmt.Sprintf("memory-dust-%v", char.Name()))

	char.AddMod(core.CharStatMod{
		Key:    "memory",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			val := make([]float64, core.EndStatType)
			atk := 0.03 + 0.01*float64(r)
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

}
