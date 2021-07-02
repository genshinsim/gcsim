package mappa

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("mappa mare", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	stacks := 0
	dur := 0
	s.AddOnReaction(func(t def.Target, ds *def.Snapshot) {
		if ds.ActorIndex == c.CharIndex() {
			stacks++
			if stacks > 2 {
				stacks = 2
				dur = s.Frame() + 600
			}
		}
	}, fmt.Sprintf("mappa-%v", c.Name()))

	dmg := 0.06 + float64(r)*0.02

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
		Key: "mappa",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			return m, dur > s.Frame()
		},
		Expiry: -1,
	})

}
