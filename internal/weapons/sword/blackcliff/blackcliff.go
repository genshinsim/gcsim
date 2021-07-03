package blackcliff

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("blackcliff longsword", weapon)
}

//After defeating an enemy, ATK is increased by 12/15/18/21/24% for 30s.
//This effect has a maximum of 3 stacks, and the duration of each stack is independent of the others.
func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {

	atk := 0.09 + float64(r)*0.03
	index := 0
	stacks := []int{-1, -1, -1}

	m := make([]float64, def.EndStatType)
	c.AddMod(def.CharStatMod{
		Key: "blackcliff",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			count := 0
			for _, v := range stacks {
				if v > s.Frame() {
					count++
				}
			}
			m[def.ATKP] = atk * float64(count)
			return m, true
		},
		Expiry: -1,
	})

	s.AddOnTargetDefeated(func(t def.Target) {
		stacks[index] = s.Frame() + 1800
		index++
		if index == 3 {
			index = 0
		}
	}, fmt.Sprintf("blackcliff-longsword-%v", c.Name()))
}
