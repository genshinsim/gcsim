package prototype

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("prototype starglitter", weapon)
}

//After using an Elemental Skill, increases Normal and Charged Attack DMG by 8% for 12s. Max 2 stacks.
func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {

	expiry := 0
	atk := 0.06 + 0.02*float64(r)
	stacks := 0
	//add on crit effect
	s.AddEventHook(func(s core.Sim) bool {
		if s.ActiveCharIndex() != c.CharIndex() {
			return false
		}
		if expiry < s.Frame() {
			stacks = 0
		}
		stacks++
		if stacks > 2 {
			stacks = 2
		}
		expiry = s.Frame() + 720
		return false
	}, fmt.Sprintf("prototype-starglitter-%v", c.Name()), core.PostSkillHook)

	val := make([]float64, core.EndStatType)
	c.AddMod(core.CharStatMod{
		Key:    "prototype",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if a != core.AttackTagNormal && a != core.AttackTagExtra {
				return nil, false
			}
			if expiry < s.Frame() {
				stacks = 0
				return nil, false
			}
			val[core.ATKP] = atk * float64(stacks)
			return val, true
		},
	})

}
