package kitain

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("kitain cross spear", weapon)
	core.RegisterWeaponFunc("kitaincrossspear", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	m := make([]float64, core.EndStatType)
	base := 0.045 + float64(r)*0.015
	regen := 2.5 + float64(r)*0.5

	m[core.PyroP] = base
	m[core.HydroP] = base
	m[core.CryoP] = base
	m[core.ElectroP] = base
	m[core.AnemoP] = base
	m[core.GeoP] = base
	m[core.EleP] = base
	m[core.PhyP] = base
	m[core.DendroP] = base

	char.AddMod(core.CharStatMod{
		Expiry: -1,
		Key:    "",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, true
		},
	})

	icd := 0
	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if ds.ActorIndex != char.CharIndex() {
			return false
		}
		if ds.AttackTag != core.AttackTagElementalArt {
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
}
