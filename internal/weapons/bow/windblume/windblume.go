package windblume

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("windblume ode", weapon)
	core.RegisterWeaponFunc("windblumeode", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	m := make([]float64, core.EndStatType)
	m[core.ATKP] = 0.12 + float64(r)*0.04

	// Effect should always apply BEFORE the skill hits
	c.Events.Subscribe(core.PreSkill, func(args ...interface{}) bool {
		char.AddMod(core.CharStatMod{
			Key: "windblume",
			Amount: func() ([]float64, bool) {
				return m, true
			},
			Expiry: c.F + 360,
		})
		return false
	}, fmt.Sprintf("windblume-%v", char.Name()))

	return "windblumeode"
}
