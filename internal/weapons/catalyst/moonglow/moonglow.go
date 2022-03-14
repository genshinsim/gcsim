package moonglow

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("everlastingmoonglow", weapon)
	core.RegisterWeaponFunc("everlasting moonglow", weapon)
	core.RegisterWeaponFunc("moonglow", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	mheal := make([]float64, core.EndStatType)
	mheal[core.Heal] = 0.075 + float64(r)*0.025
	char.AddMod(coretype.CharStatMod{
		Key: "moonglow-heal-bonus",
		Amount: func() ([]float64, bool) {
			return mheal, true
		},
		Expiry: -1,
	})

	nabuff := 0.0005 + float64(r)*0.0005
	matk := make([]float64, core.EndStatType)
	char.AddPreDamageMod(coretype.PreDamageMod{
		Key: "moonglow-na-bonus",
		Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
			if atk.Info.AttackTag != coretype.AttackTagNormal {
				return nil, false
			}

			matk[core.ATK] = nabuff * char.MaxHP()
			return matk, true
		},
		Expiry: -1,
	})

	icd, dur := -1, -1
	c.Subscribe(core.PreBurst, func(args ...interface{}) bool {
		if c.ActiveChar != char.Index() {
			return false
		}
		dur = c.Frame + 720 // 12s

		return false
	}, fmt.Sprintf("moonglow-onburst-%v", char.Name()))

	c.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*coretype.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return false
		}
		if atk.Info.AttackTag != coretype.AttackTagNormal {
			return false
		}
		if dur < c.Frame || icd > c.Frame {
			return false
		}

		char.AddEnergy("moonglow", 0.6)
		icd = c.Frame + 6 // 0.1s

		return false
	}, fmt.Sprintf("moonglow-energy-%v", char.Name()))

	return "everlastingmoonglow"
}
