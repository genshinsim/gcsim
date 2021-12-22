package blackcliff

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("blackcliff warbow", weapon)
	core.RegisterWeaponFunc("blackcliff slasher", weapon)
	core.RegisterWeaponFunc("blackcliff agate", weapon)
	core.RegisterWeaponFunc("blackcliff pole", weapon)
	core.RegisterWeaponFunc("blackcliff longsword", weapon)
	core.RegisterWeaponFunc("blackcliffagate", weapon)
	core.RegisterWeaponFunc("blackclifflongsword", weapon)
	core.RegisterWeaponFunc("blackcliffpole", weapon)
	core.RegisterWeaponFunc("blackcliffslasher", weapon)
	core.RegisterWeaponFunc("blackcliffwarbow", weapon)

}

//After defeating an enemy, ATK is increased by 12/15/18/21/24% for 30s.
//This effect has a maximum of 3 stacks, and the duration of each stack is independent of the others.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) {

	atk := 0.09 + float64(r)*0.03
	index := 0
	stacks := []int{-1, -1, -1}

	m := make([]float64, core.EndStatType)
	char.AddMod(core.CharStatMod{
		Key: "blackcliff",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			count := 0
			for _, v := range stacks {
				if v > c.F {
					count++
				}
			}
			m[core.ATKP] = atk * float64(count)
			return m, true
		},
		Expiry: -1,
	})

	c.Events.Subscribe(core.OnTargetDied, func(args ...interface{}) bool {
		stacks[index] = c.F + 1800
		index++
		if index == 3 {
			index = 0
		}
		return false
	}, fmt.Sprintf("blackcliff-%v", char.Name()))

}
