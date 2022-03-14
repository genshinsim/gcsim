package whiteblind

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("whiteblind", weapon)
}

//On hit, Normal or Charged Attacks increase ATK and DEF by 6/7.5/9/10.5/12% for 6s.
//Max 4 stacks (24/30/36/42/48% total). This effect can only occur once every 0.5s.
func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	stacks := 0
	icd := 0
	duration := 0

	c.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*coretype.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return false
		}
		if atk.Info.AttackTag != coretype.AttackTagNormal && atk.Info.AttackTag != coretype.AttackTagExtra {
			return false
		}
		if icd > c.Frame {
			return false
		}
		if duration < c.Frame {
			stacks = 0
		}
		duration = c.Frame + 360
		icd = c.Frame + 30
		stacks++
		if stacks > 4 {
			stacks = 4
		}
		return false
	}, fmt.Sprintf("whiteblind-%v", char.Name()))

	amt := 0.045 + float64(r)*0.015

	val := make([]float64, core.EndStatType)
	char.AddMod(coretype.CharStatMod{
		Key:    "whiteblind",
		Expiry: -1,
		Amount: func() ([]float64, bool) {
			if duration < c.Frame {
				stacks = 0
				return nil, false
			}
			val[core.ATKP] = amt * float64(stacks)
			val[core.DEFP] = amt * float64(stacks)
			return val, true
		},
	})
	return "whiteblind"
}
