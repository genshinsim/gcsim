package kitain

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("kitain cross spear", weapon)
	core.RegisterWeaponFunc("kitaincrossspear", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	m := make([]float64, core.EndStatType)
	base := 0.045 + float64(r)*0.015
	regen := 2.5 + float64(r)*0.5

	m[core.DmgP] = base

	char.AddMod(core.CharStatMod{
		Expiry: -1,
		Key:    "kitain-skill-dmg-buff",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if a == core.AttackTagElementalArt || a == core.AttackTagElementalArtHold {
				return m, true
			}
			return nil, false
		},
	})

	icd := 0
	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if atk.Info.AttackTag != core.AttackTagElementalArt {
			return false
		}
		if icd > c.F {
			return false
		}
		icd = c.F + 600 //once every 10 seconds
		char.AddEnergy(-3)
		for i := 120; i <= 360; i += 120 {
			char.AddTask(func() {
				char.AddEnergy(regen)
			}, "kitain-restore", i)
		}
		return false
	}, fmt.Sprintf("kitain-%v", char.Name()))
	return "kitaincrossspear"
}
