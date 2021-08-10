package skyward

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("skyward atlas", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	dmg := 0.09 + float64(r)*0.03
	atk := 1.2 + float64(r)*0.4

	icd := 0

	s.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag != core.AttackTagNormal {
			return
		}
		if icd > s.Frame() {
			return
		}
		if s.Rand().Float64() < 0.5 {
			return
		}
		d := c.Snapshot(
			"Skyward Atlas Proc",
			core.AttackTagWeaponSkill,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Physical,
			100,
			atk,
		)
		c.QueueDmg(&d, 1)
		for i := 0; i < 6; i++ {
			x := d.Clone()
			c.QueueDmg(&x, i*150)
		}
		icd = s.Frame() + 1800
	}, fmt.Sprintf("skyward-atlast-%v", c.Name()))

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
	c.AddMod(core.CharStatMod{
		Key:    "skyward-atlast",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, true
		},
	})
}
