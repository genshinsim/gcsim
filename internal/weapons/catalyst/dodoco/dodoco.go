package dodoco

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("dodoco tales", weapon)
	core.RegisterWeaponFunc("dodocotales", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	atkExpiry := 0
	dmgExpiry := 0

	m := make([]float64, core.EndStatType)
	m[core.DmgP] = .12 + float64(r)*.04
	char.AddPreDamageMod(coretype.PreDamageMod{
		Key: "dodoco-ca",
		Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
			if atk.Info.AttackTag != coretype.AttackTagExtra {
				return nil, false
			}
			return m, dmgExpiry > c.Frame
		},
		Expiry: -1,
	})

	n := make([]float64, core.EndStatType)
	n[core.ATKP] = .06 + float64(r)*0.02
	char.AddMod(coretype.CharStatMod{
		Key: "dodoco atk",
		Amount: func() ([]float64, bool) {
			return n, atkExpiry > c.Frame
		},
		Expiry: -1,
	})

	c.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*coretype.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return false
		}
		switch atk.Info.AttackTag {
		case coretype.AttackTagNormal:
			dmgExpiry = c.Frame + 360
		case coretype.AttackTagExtra:
			atkExpiry = c.Frame + 360
		}
		return false
	}, fmt.Sprintf("dodoco-%v", char.Name()))

	return "dodocotales"
}
