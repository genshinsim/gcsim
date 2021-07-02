package heartofdepth

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterSetFunc("heart of depth", New)
}

func New(c def.Character, s def.Sim, log def.Logger, count int) {
	if count >= 2 {
		m := make([]float64, def.EndStatType)
		m[def.HydroP] = 0.15
		c.AddMod(def.CharStatMod{
			Key: "hod-2pc",
			Amount: func(a def.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		key := fmt.Sprintf("%v-hod-4pc", c.Name())

		s.AddEventHook(func(s def.Sim) bool {
			s.AddStatus(key, 900) //activate for 15 seoncds
			return false
		}, fmt.Sprintf("hod4-%v", c.Name()), def.PostSkillHook)

		m := make([]float64, def.EndStatType)
		m[def.DmgP] = 0.3
		c.AddMod(def.CharStatMod{
			Key: "hod-4pc",
			Amount: func(ds def.AttackTag) ([]float64, bool) {
				if s.Status(key) == 0 {
					return nil, false
				}
				if ds != def.AttackTagNormal && ds != def.AttackTagExtra {
					return nil, false
				}
				return m, true
			},
			Expiry: -1,
		})
	}
	//add flat stat to char
}
