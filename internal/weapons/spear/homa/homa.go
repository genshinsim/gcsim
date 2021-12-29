package homa

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("staff of homa", weapon)
	core.RegisterWeaponFunc("staffofhoma", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	//add on hit effect to sim?
	m := make([]float64, core.EndStatType)
	m[core.HPP] = 0.15 + float64(r)*0.05
	atkp := 0.006 + float64(r)*0.002
	lowhp := 0.008 + float64(r)*0.002

	char.AddMod(core.CharStatMod{
		Key: "homa hp bonus",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			per := atkp
			if char.HP()/char.MaxHP() <= 0.5 {
				per += lowhp
			}
			// c.Log.Debugw("homa bonus atk%", "frame", c.F, "char", char.CharIndex(), "event", core.LogSnapshotEvent, "max-hp", char.MaxHP(), "percent", char.HP()/char.MaxHP(), "per", per)
			m[core.ATK] = per * char.MaxHP()
			return m, true
		},
		Expiry: -1,
	})
	return "staffofhoma"
}
