package heartofdepth

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterSetFunc("heart of depth", New)
}

func New(c core.Character, s core.Sim, log core.Logger, count int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.HydroP] = 0.15
		c.AddMod(core.CharStatMod{
			Key: "hod-2pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		key := fmt.Sprintf("%v-hod-4pc", c.Name())

		s.AddEventHook(func(s core.Sim) bool {
			s.AddStatus(key, 900) //activate for 15 seoncds
			return false
		}, fmt.Sprintf("hod4-%v", c.Name()), core.PostSkillHook)

		m := make([]float64, core.EndStatType)
		m[core.DmgP] = 0.3
		c.AddMod(core.CharStatMod{
			Key: "hod-4pc",
			Amount: func(ds core.AttackTag) ([]float64, bool) {
				if s.Status(key) == 0 {
					return nil, false
				}
				if ds != core.AttackTagNormal && ds != core.AttackTagExtra {
					return nil, false
				}
				return m, true
			},
			Expiry: -1,
		})
	}
	//add flat stat to char
}
