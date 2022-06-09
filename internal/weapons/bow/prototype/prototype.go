package generic

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("prototype crescent", weapon)
	core.RegisterWeaponFunc("prototypecrescent", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {

	m := make([]float64, core.EndStatType)
	m[core.ATKP] = 0.27 + float64(r)*0.09

	//add on hit effect
	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if !atk.Info.HitWeakPoint {
			return false
		}
		char.AddMod(core.CharStatMod{
			Key:    "prototype-crescent",
			Expiry: c.F + 60*10,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
		return false
	}, fmt.Sprintf("prototype-crescent-%v", char.Name()))

	return "prototypecrescent"
}
