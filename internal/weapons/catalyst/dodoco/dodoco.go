package dodoco

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("dodoco tales", weapon)
	core.RegisterWeaponFunc("dodocotales", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	atkExpiry := 0
	dmgExpiry := 0

	m := make([]float64, core.EndStatType)
	m[core.DmgP] = .12 + float64(r)*.04
	char.AddMod(core.CharStatMod{
		Key: "dodoco ca",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if a != core.AttackTagExtra {
				return nil, false
			}
			return m, dmgExpiry > c.F
		},
		Expiry: -1,
	})

	n := make([]float64, core.EndStatType)
	n[core.ATKP] = .06 + float64(r)*0.02
	char.AddMod(core.CharStatMod{
		Key: "dodoco atk",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return n, atkExpiry > c.F
		},
		Expiry: -1,
	})

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		switch atk.Info.AttackTag {
		case core.AttackTagNormal:
			dmgExpiry = c.F + 360
		case core.AttackTagExtra:
			atkExpiry = c.F + 360
		}
		return false
	}, fmt.Sprintf("dodoco-%v", char.Name()))

	return "dodocotales"
}
