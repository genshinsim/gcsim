package hamayumi

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("iram", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {

	//28,56,96,96
	stacks := 0
	amt := []float64{0, .28, .56, .96, .96}
	expiry := 0

	char.AddMod(core.CharStatMod{
		Key:    "iram",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			val := make([]float64, core.EndStatType)
			if expiry < c.F {
				stacks = 0
				return nil, false
			}
			val[core.ATKP] = amt[stacks]
			return val, true
		},
	})

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		switch atk.Info.AttackTag {
		case core.AttackTagNormal:
		case core.AttackTagExtra:
		case core.AttackTagElementalArt:
		case core.AttackTagElementalBurst:
		default:
			return false
		}

		if stacks < 4 {
			stacks++
		}
		expiry = c.F + 600

		return false
	}, fmt.Sprintf("iram-%v", char.Name()))
}
