package windblume

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("windblume ode", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {

	dur := 0
	//add on hit effect
	s.AddEventHook(func(s def.Sim) bool {
		dur = s.Frame() + 360
		return false
	}, fmt.Sprintf("windblume-%v", c.Name()), def.PostSkillHook)

	m := make([]float64, def.EndStatType)
	m[def.ATKP] = 0.12 + float64(r)*0.04
	c.AddMod(def.CharStatMod{
		Key: "windblume",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			if dur < s.Frame() {
				return nil, false
			}
			return m, true
		},
		Expiry: -1,
	})
}
