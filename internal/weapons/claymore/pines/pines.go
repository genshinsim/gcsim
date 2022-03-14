package pines

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("song of broken pines", weapon)
	core.RegisterWeaponFunc("songofbrokenpines", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	m := make([]float64, core.EndStatType)
	m[core.ATKP] = 0.12 + float64(r)*0.04
	char.AddMod(coretype.CharStatMod{
		Key: "pines-atk",
		Amount: func() ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	val := make([]float64, core.EndStatType)
	val[core.ATKP] = 0.15 + 0.05*float64(r)
	val[core.AtkSpd] = 0.09 + 0.03*float64(r)

	icd := 0
	stacks := 0
	cooldown := 0

	c.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*coretype.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return false
		}
		if atk.Info.AttackTag != coretype.AttackTagNormal && atk.Info.AttackTag != coretype.AttackTagExtra {
			return false
		}
		if cooldown > c.Frame {
			return false
		}
		if icd > c.Frame {
			return false
		}
		icd = c.Frame + 12
		stacks++
		if stacks == 4 {
			stacks = 0
			c.AddStatus("pines", 720)
			cooldown = c.Frame + 1200
			for _, char := range c.Chars {
				char.AddMod(coretype.CharStatMod{
					Key: "pines-proc",
					Amount: func() ([]float64, bool) {
						return val, true
					},
					Expiry: c.Frame + 720,
				})
			}
		}
		return false
	}, fmt.Sprintf("pines-%v", char.Name()))

	return "songofbrokenpines"
}
