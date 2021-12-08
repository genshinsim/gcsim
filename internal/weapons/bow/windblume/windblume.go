package windblume

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("windblume ode", weapon)
	core.RegisterWeaponFunc("windblumeode", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {

	dur := 0
	//add on hit effect
	c.Events.Subscribe(core.PostSkill, func(args ...interface{}) bool {
		dur = c.F + 360
		return false
	}, fmt.Sprintf("windblume-%v", char.Name()))

	var m [core.EndStatType]float64
	m[core.ATKP] = 0.12 + float64(r)*0.04
	char.AddMod(core.CharStatMod{
		Key: "windblume",
		Amount: func(a core.AttackTag) ([core.EndStatType]float64, bool) {
			if dur < c.F {
				return m, false
			}
			return m, true
		},
		Expiry: -1,
	})
}
