package whiteblind

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("whiteblind", weapon)
}

//On hit, Normal or Charged Attacks increase ATK and DEF by 6/7.5/9/10.5/12% for 6s.
//Max 4 stacks (24/30/36/42/48% total). This effect can only occur once every 0.5s.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	stacks := 0
	icd := 0
	duration := 0

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
		if duration < c.F {
			stacks = 0
		}
		duration = c.F + 360
		icd = c.F + 30
		stacks++
		if stacks > 4 {
			stacks = 4
		}
		return false
	}, fmt.Sprintf("whiteblind-%v", char.Name()))

	amt := 0.045 + float64(r)*0.015

	val := make([]float64, core.EndStatType)
	char.AddMod(core.CharStatMod{
		Key:    "whiteblind",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if duration < c.F {
				stacks = 0
				return nil, false
			}
			val[core.ATKP] = amt * float64(stacks)
			val[core.DEFP] = amt * float64(stacks)
			return val, true
		},
	})
}
