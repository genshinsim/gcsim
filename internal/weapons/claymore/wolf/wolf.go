package generic

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("wolf's gravestone", weapon)
	core.RegisterWeaponFunc("wolfsgravestone", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {

	val := make([]float64, core.EndStatType)
	val[core.ATKP] = 0.15 + 0.05*float64(r)
	char.AddMod(core.CharStatMod{
		Key:    "wolf-flat",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return val, true
		},
	})

	bonus := make([]float64, core.EndStatType)
	bonus[core.ATKP] = 0.3 + 0.1*float64(r)
	icd := 0

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		if !c.Flags.DamageMode {
			return false //ignore as we not tracking HP
		}

		atk := args[1].(*core.AttackEvent)
		t := args[0].(core.Target)
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if icd > c.F {
			return false
		}

		if t.HP()/t.MaxHP() > 0.3 {
			return false
		}
		icd = c.F + 1800 //every 30 seconds

		for _, char := range c.Chars {
			char.AddMod(core.CharStatMod{
				Key:    "wolf-proc",
				Expiry: c.F + 720,
				Amount: func(a core.AttackTag) ([]float64, bool) {
					return bonus, true
				},
			})
		}
		return false
	}, fmt.Sprintf("wolf-%v", char.Name()))
	return "wolfsgravestone"
}
