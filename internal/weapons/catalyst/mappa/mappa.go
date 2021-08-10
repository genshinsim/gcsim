package mappa

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("mappa mare", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	stacks := 0
	dur := 0
	s.AddOnReaction(func(t core.Target, ds *core.Snapshot) {
		if ds.ActorIndex == c.CharIndex() {
			stacks++
			if stacks > 2 {
				stacks = 2
				dur = s.Frame() + 600
			}
		}
	}, fmt.Sprintf("mappa-%v", c.Name()))

	dmg := 0.06 + float64(r)*0.02

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
		Key: "mappa",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, dur > s.Frame()
		},
		Expiry: -1,
	})

}
