package generic

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("prototype crescent", weapon)
	core.RegisterWeaponFunc("prototypecrescent", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {

	dur := 0
	key := fmt.Sprintf("prototype-crescent-%v", char.Name())
	//add on hit effect
	c.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*coretype.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return false
		}
		if atk.Info.HitWeakPoint {
			dur = c.Frame + 600
		}
		return false
	}, key)

	m := make([]float64, core.EndStatType)
	m[core.ATKP] = 0.27 + float64(r)*0.09
	char.AddMod(coretype.CharStatMod{
		Key: "prototype-crescent",
		Amount: func() ([]float64, bool) {
			if dur < c.Frame {
				return nil, false
			}
			return m, true
		},
		Expiry: -1,
	})

	return "prototypecrescent"
}
