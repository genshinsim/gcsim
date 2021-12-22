package pines

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("song of broken pines", weapon)
	core.RegisterWeaponFunc("songofbrokenpines", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	m := make([]float64, core.EndStatType)
	m[core.ATKP] = 0.12 + float64(r)*0.04
	char.AddMod(core.CharStatMod{
		Key: "pines-atk",
		Amount: func(a core.AttackTag) ([]float64, bool) {
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

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if atk.Info.AttackTag != core.AttackTagNormal && atk.Info.AttackTag != core.AttackTagExtra {
			return false
		}
		if cooldown > c.F {
			return false
		}
		if icd > c.F {
			return false
		}
		icd = c.F + 12
		stacks++
		if stacks == 4 {
			stacks = 0
			c.Status.AddStatus("pines", 720)
			cooldown = c.F + 1200
			for _, char := range c.Chars {
				char.AddMod(core.CharStatMod{
					Key: "pines-proc",
					Amount: func(a core.AttackTag) ([]float64, bool) {
						return val, true
					},
					Expiry: c.F + 720,
				})
			}
		}
		return false
	}, fmt.Sprintf("pines-%v", char.Name()))
}
