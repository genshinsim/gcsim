package generic

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("wolf's gravestone", weapon)
	core.RegisterWeaponFunc("wolfsgravestone", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {

	val := make([]float64, core.EndStatType)
	val[core.ATKP] = 0.15 + 0.05*float64(r)
	char.AddMod(coretype.CharStatMod{
		Key:    "wolf-flat",
		Expiry: -1,
		Amount: func() ([]float64, bool) {
			return val, true
		},
	})

	bonus := make([]float64, core.EndStatType)
	bonus[core.ATKP] = 0.3 + 0.1*float64(r)
	icd := 0

	c.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		if !c.Flags.DamageMode {
			return false //ignore as we not tracking HP
		}

		atk := args[1].(*coretype.AttackEvent)
		t := args[0].(coretype.Target)
		if atk.Info.ActorIndex != char.Index() {
			return false
		}
		if icd > c.Frame {
			return false
		}

		if t.HP()/t.MaxHP() > 0.3 {
			return false
		}
		icd = c.Frame + 1800 //every 30 seconds

		for _, char := range c.Chars {
			char.AddMod(coretype.CharStatMod{
				Key:    "wolf-proc",
				Expiry: c.Frame + 720,
				Amount: func() ([]float64, bool) {
					return bonus, true
				},
			})
		}
		return false
	}, fmt.Sprintf("wolf-%v", char.Name()))
	return "wolfsgravestone"
}
