package skyward

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("skyward atlas", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	dmg := 0.09 + float64(r)*0.03
	atk := 1.2 + float64(r)*0.4

	icd := 0

	s.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag != def.AttackTagNormal {
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
			def.AttackTagWeaponSkill,
			def.ICDTagNone,
			def.ICDGroupDefault,
			def.StrikeTypeDefault,
			def.Physical,
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

	m := make([]float64, def.EndStatType)
	m[def.PyroP] = dmg
	m[def.HydroP] = dmg
	m[def.CryoP] = dmg
	m[def.ElectroP] = dmg
	m[def.AnemoP] = dmg
	m[def.GeoP] = dmg
	m[def.EleP] = dmg
	m[def.PhyP] = dmg
	m[def.DendroP] = dmg
	c.AddMod(def.CharStatMod{
		Key:    "skyward-atlast",
		Expiry: -1,
		Amount: func(a def.AttackTag) ([]float64, bool) {
			return m, true
		},
	})
}
