package prototype

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("prototype starglitter", weapon)
	core.RegisterWeaponFunc("prototypestarglitter", weapon)
}

//After using an Elemental Skill, increases Normal and Charged Attack DMG by 8% for 12s. Max 2 stacks.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {

	expiry := 0
	atk := 0.06 + 0.02*float64(r)
	stacks := 0
	//add on crit effect
	c.Events.Subscribe(core.PostSkill, func(args ...interface{}) bool {
		if c.ActiveChar != char.CharIndex() {
			return false
		}
		if expiry < c.F {
			stacks = 0
		}
		stacks++
		if stacks > 2 {
			stacks = 2
		}
		expiry = c.F + 720
		return false
	}, fmt.Sprintf("prototype-starglitter-%v", char.Name()))

	char.AddMod(core.CharStatMod{
		Key:    "prototype",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			val := make([]float64, core.EndStatType)
			if a != core.AttackTagNormal && a != core.AttackTagExtra {
				return nil, false
			}
			if expiry < c.F {
				stacks = 0
				return nil, false
			}
			val[core.ATKP] = atk * float64(stacks)
			return val, true
		},
	})
	return "prototypestarglitter"
}
