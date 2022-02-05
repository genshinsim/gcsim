package moonglow

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("everlastingmoonglow", weapon)
	core.RegisterWeaponFunc("everlasting moonglow", weapon)
	core.RegisterWeaponFunc("moonglow", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	mheal := make([]float64, core.EndStatType)
	mheal[core.Heal] = 0.075 + float64(r)*0.025
	char.AddMod(core.CharStatMod{
		Key: "moonglow-heal-bonus",
		Amount: func() ([]float64, bool) {
			return mheal, true
		},
		Expiry: -1,
	})

	nabuff := 0.0005 + float64(r)*0.0005
	matk := make([]float64, core.EndStatType)
	char.AddPreDamageMod(core.PreDamageMod{
		Key: "moonglow-na-bonus",
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			if atk.Info.AttackTag != core.AttackTagNormal {
				return nil, false
			}

			matk[core.ATK] = nabuff * char.MaxHP()
			return matk, true
		},
		Expiry: -1,
	})

	icd, dur := -1, -1
	c.Events.Subscribe(core.PreBurst, func(args ...interface{}) bool {
		if c.ActiveChar != char.CharIndex() {
			return false
		}
		dur = c.F + 720 // 12s

		return false
	}, fmt.Sprintf("moonglow-onburst-%v", char.Name()))

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if atk.Info.AttackTag != core.AttackTagNormal {
			return false
		}
		if dur < c.F || icd > c.F {
			return false
		}

		char.AddEnergy(0.6)
		icd = c.F + 6 // 0.1s

		return false
	}, fmt.Sprintf("moonglow-energy-%v", char.Name()))

	return "everlastingmoonglow"
}
