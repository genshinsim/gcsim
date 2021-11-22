package mappa

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("mappa mare", weapon)
	core.RegisterWeaponFunc("mappamare", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	stacks := 0
	dur := 0
	c.Events.Subscribe(core.OnReactionOccured, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if ds.ActorIndex == char.CharIndex() {
			stacks++
			if stacks > 2 {
				stacks = 2
				dur = c.F + 600
			}
		}
		return false
	}, fmt.Sprintf("mappa-%v", char.Name()))

	dmg := 0.06 + float64(r)*0.02

	m := make([]float64, core.EndStatType)

	m[core.PyroP] = dmg
	m[core.HydroP] = dmg
	m[core.CryoP] = dmg
	m[core.ElectroP] = dmg
	m[core.AnemoP] = dmg
	m[core.GeoP] = dmg
	m[core.DendroP] = dmg

	char.AddMod(core.CharStatMod{
		Key: "mappa",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, dur > c.F
		},
		Expiry: -1,
	})

}
