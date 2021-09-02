package hamayumi

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("iram", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {

	//28,56,96,96
	stacks := 0
	amt := []float64{0, .28, .56, .96, .96}
	expiry := 0

	val := make([]float64, core.EndStatType)
	char.AddMod(core.CharStatMod{
		Key:    "iram",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if expiry < c.F {
				stacks = 0
				return nil, false
			}
			val[core.ATKP] = amt[stacks]
			return val, true
		},
	})

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if ds.ActorIndex != char.CharIndex() {
			return false
		}
		switch ds.AttackTag {
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
