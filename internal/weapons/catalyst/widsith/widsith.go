package widsith

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("the widsith", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	last := 0
	expiry := 0

	atk := .45 + float64(r)*0.15
	em := 180 + float64(r)*60
	dmg := .36 + float64(r)*.12

	m := make([]float64, def.EndStatType)

	c.AddMod(def.CharStatMod{
		Key: "widsith",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			return m, expiry > s.Frame()
		},
		Expiry: -1,
	})

	s.AddEventHook(func(s def.Sim) bool {
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
			m[def.EM] = em
			m[def.PyroP] = 0
			m[def.HydroP] = 0
			m[def.CryoP] = 0
			m[def.ElectroP] = 0
			m[def.AnemoP] = 0
			m[def.GeoP] = 0
			m[def.EleP] = 0
			m[def.PhyP] = 0
			m[def.DendroP] = 0
			m[def.ATKP] = 0
			log.Debugw("widsith proc'd", "frame", s.Frame(), "event", def.LogWeaponEvent, "char", c.CharIndex(), "stat", "em", "expiring", expiry)
		case 1:
			m[def.EM] = 0
			m[def.PyroP] = dmg
			m[def.HydroP] = dmg
			m[def.CryoP] = dmg
			m[def.ElectroP] = dmg
			m[def.AnemoP] = dmg
			m[def.GeoP] = dmg
			m[def.EleP] = dmg
			m[def.PhyP] = dmg
			m[def.DendroP] = dmg
			m[def.ATKP] = 0
			log.Debugw("widsith proc'd", "frame", s.Frame(), "event", def.LogWeaponEvent, "char", c.CharIndex(), "stat", "dmg%", "expiring", expiry)
		default:
			m[def.EM] = 0
			m[def.PyroP] = 0
			m[def.HydroP] = 0
			m[def.CryoP] = 0
			m[def.ElectroP] = 0
			m[def.AnemoP] = 0
			m[def.GeoP] = 0
			m[def.EleP] = 0
			m[def.PhyP] = 0
			m[def.DendroP] = 0
			m[def.ATKP] = atk
			log.Debugw("widsith proc'd", "frame", s.Frame(), "event", def.LogWeaponEvent, "char", c.CharIndex(), "stat", "atk%", "expiring", expiry)
		}

		return false
	}, fmt.Sprintf("width-%v", c.Name()), def.PostSwapHook)
}
