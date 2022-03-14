package prototype

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("prototype starglitter", weapon)
	core.RegisterWeaponFunc("prototypestarglitter", weapon)
}

//After using an Elemental Skill, increases Normal and Charged Attack DMG by 8% for 12s. Max 2 stacks.
func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {

	expiry := 0
	atkbonus := 0.06 + 0.02*float64(r)
	stacks := 0
	//add on crit effect
	c.Subscribe(core.PreSkill, func(args ...interface{}) bool {
		if c.ActiveChar != char.Index() {
			return false
		}
		if expiry < c.Frame {
			stacks = 0
		}
		stacks++
		if stacks > 2 {
			stacks = 2
		}
		expiry = c.Frame + 720
		return false
	}, fmt.Sprintf("prototype-starglitter-%v", char.Name()))

	char.AddPreDamageMod(coretype.PreDamageMod{
		Key:    "prototype",
		Expiry: -1,
		Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
			val := make([]float64, core.EndStatType)
			if atk.Info.AttackTag != coretype.AttackTagNormal && atk.Info.AttackTag != coretype.AttackTagExtra {
				return nil, false
			}
			if expiry < c.Frame {
				stacks = 0
				return nil, false
			}
			val[core.ATKP] = atkbonus * float64(stacks)
			return val, true
		},
	})
	return "prototypestarglitter"
}
