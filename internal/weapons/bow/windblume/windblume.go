package windblume

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("windblume ode", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {

	dur := 0
	//add on hit effect
	s.AddEventHook(func(s core.Sim) bool {
		dur = s.Frame() + 360
		return false
	}, fmt.Sprintf("windblume-%v", c.Name()), core.PostSkillHook)

	m := make([]float64, core.EndStatType)
	m[core.ATKP] = 0.12 + float64(r)*0.04
	c.AddMod(core.CharStatMod{
		Key: "windblume",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if dur < s.Frame() {
				return nil, false
			}
			return m, true
		},
		Expiry: -1,
	})
}
