package darkironsword

import (
	"fmt"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("darkironsword", weapon)
}

// Overloaded
// Upon causing an Overloaded, Superconduct, Electro-Charged, or an Electro-infused Swirl reaction,
// ATK is increased by 20/25/30/35/40% for 12s.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	dur := 12 * 60
	m := make([]float64, core.EndStatType)
	m[core.ATKP] = 0.15 + float64(r)*0.05

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)

		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}

		//ignore if character not on field
		if c.ActiveChar != char.CharIndex() {
			return false
		}

		switch atk.Info.AttackTag {
		case core.AttackTagSuperconductDamage,
			core.AttackTagECDamage,
			core.AttackTagOverloadDamage,
			core.AttackTagSwirlElectro:
			char.AddMod(core.CharStatMod{
				Key: "darkironsword",
				Amount: func() ([]float64, bool) {
					return m, true
				},
				Expiry: c.F + dur,
			})
		}

		return false
	}, fmt.Sprintf("darkironsword-%v", char.Name()))
	return "darkironsword"
}
