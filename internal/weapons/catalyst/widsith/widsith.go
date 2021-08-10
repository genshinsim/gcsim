package widsith

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("the widsith", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	last := 0
	expiry := 0

	atk := .45 + float64(r)*0.15
	em := 180 + float64(r)*60
	dmg := .36 + float64(r)*.12

	m := make([]float64, core.EndStatType)

	c.AddMod(core.CharStatMod{
		Key: "widsith",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, expiry > s.Frame()
		},
		Expiry: -1,
	})

	s.AddEventHook(func(s core.Sim) bool {
		//ignore if char is not the active one
		if s.ActiveCharIndex() != c.CharIndex() {
			return false
		}
		//if char is the active one then we just came on to field
		if last != 0 && s.Frame()-last < 1800 { //30 sec icd
			return false
		}
		last = s.Frame()
		expiry = s.Frame() + 600 //10 sec duration
		//random 1 of 3
		i := s.Rand().Intn(3)

		switch i {
		case 0:
			m[core.EM] = em
			m[core.PyroP] = 0
			m[core.HydroP] = 0
			m[core.CryoP] = 0
			m[core.ElectroP] = 0
			m[core.AnemoP] = 0
			m[core.GeoP] = 0
			m[core.EleP] = 0
			m[core.PhyP] = 0
			m[core.DendroP] = 0
			m[core.ATKP] = 0
			log.Debugw("widsith proc'd", "frame", s.Frame(), "event", core.LogWeaponEvent, "char", c.CharIndex(), "stat", "em", "expiring", expiry)
		case 1:
			m[core.EM] = 0
			m[core.PyroP] = dmg
			m[core.HydroP] = dmg
			m[core.CryoP] = dmg
			m[core.ElectroP] = dmg
			m[core.AnemoP] = dmg
			m[core.GeoP] = dmg
			m[core.EleP] = dmg
			m[core.PhyP] = dmg
			m[core.DendroP] = dmg
			m[core.ATKP] = 0
			log.Debugw("widsith proc'd", "frame", s.Frame(), "event", core.LogWeaponEvent, "char", c.CharIndex(), "stat", "dmg%", "expiring", expiry)
		default:
			m[core.EM] = 0
			m[core.PyroP] = 0
			m[core.HydroP] = 0
			m[core.CryoP] = 0
			m[core.ElectroP] = 0
			m[core.AnemoP] = 0
			m[core.GeoP] = 0
			m[core.EleP] = 0
			m[core.PhyP] = 0
			m[core.DendroP] = 0
			m[core.ATKP] = atk
			log.Debugw("widsith proc'd", "frame", s.Frame(), "event", core.LogWeaponEvent, "char", c.CharIndex(), "stat", "atk%", "expiring", expiry)
		}

		return false
	}, fmt.Sprintf("width-%v", c.Name()), core.PostSwapHook)
}
