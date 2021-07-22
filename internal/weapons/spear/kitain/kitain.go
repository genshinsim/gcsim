package kitain

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("kitain cross spear", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	m := make([]float64, def.EndStatType)
	base := 0.045 + float64(r)*0.015
	regen := 2.5 + float64(r)*0.5

	m[def.PyroP] = base
	m[def.HydroP] = base
	m[def.CryoP] = base
	m[def.ElectroP] = base
	m[def.AnemoP] = base
	m[def.GeoP] = base
	m[def.EleP] = base
	m[def.PhyP] = base
	m[def.DendroP] = base

	c.AddMod(def.CharStatMod{
		Expiry: -1,
		Key:    "",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			return m, true
		},
	})

	icd := 0
	s.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag != def.AttackTagElementalArt {
			return
		}
		if icd > s.Frame() {
			return
		}
		icd = s.Frame() + 600 //once every 10 seconds
		c.AddEnergy(-3)
		for i := 120; i <= 360; i += 120 {
			c.AddTask(func() {
				c.AddEnergy(regen)
			}, "kitain-restore", i)
		}

	}, fmt.Sprintf("kitain-%v", c.Name()))
}
