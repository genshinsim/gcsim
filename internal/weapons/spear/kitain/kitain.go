package kitain

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("kitain cross spear", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
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

	c.AddMod(core.CharStatMod{
		Expiry: -1,
		Key:    "",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, true
		},
	})

	icd := 0
	s.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag != core.AttackTagElementalArt {
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
