package skyward

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("skyward atlas", weapon)
	core.RegisterWeaponFunc("skywardatlas", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	dmg := 0.09 + float64(r)*0.03
	atk := 1.2 + float64(r)*0.4

	icd := 0

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if ds.ActorIndex != char.CharIndex() {
			return false
		}
		if ds.AttackTag != core.AttackTagNormal {
			return false
		}
		if icd > c.F {
			return false
		}
		if c.Rand.Float64() < 0.5 {
			return false
		}
		d := char.Snapshot(
			"Skyward Atlas Proc",
			core.AttackTagWeaponSkill,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Physical,
			100,
			atk,
		)
		char.QueueDmg(&d, 1)
		for i := 0; i < 6; i++ {
			x := d.Clone()
			char.QueueDmg(&x, i*150)
		}
		icd = c.F + 1800
		return false
	}, fmt.Sprintf("skyward-atlast-%v", char.Name()))

	m := make([]float64, core.EndStatType)
	m[core.PyroP] = dmg
	m[core.HydroP] = dmg
	m[core.CryoP] = dmg
	m[core.ElectroP] = dmg
	m[core.AnemoP] = dmg
	m[core.GeoP] = dmg
	m[core.EleP] = dmg
	m[core.PhyP] = dmg
	m[core.DendroP] = dmg
	char.AddMod(core.CharStatMod{
		Key:    "skyward-atlast",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, true
		},
	})
}
