package prototype

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("prototype starglitter", weapon)
}

//After using an Elemental Skill, increases Normal and Charged Attack DMG by 8% for 12s. Max 2 stacks.
func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {

	expiry := 0
	atk := 0.06 + 0.02*float64(r)
	stacks := 0
	//add on crit effect
	s.AddEventHook(func(s def.Sim) bool {
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
	}, fmt.Sprintf("prototype-starglitter-%v", c.Name()), def.PostSkillHook)

	val := make([]float64, def.EndStatType)
	c.AddMod(def.CharStatMod{
		Key:    "prototype",
		Expiry: -1,
		Amount: func(a def.AttackTag) ([]float64, bool) {
			if a != def.AttackTagNormal && a != def.AttackTagExtra {
				return nil, false
			}
			if expiry < s.Frame() {
				stacks = 0
				return nil, false
			}
			val[def.ATKP] = atk * float64(stacks)
			return val, true
		},
	})

}
