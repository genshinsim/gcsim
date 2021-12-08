package heartofdepth

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("heart of depth", New)
	core.RegisterSetFunc("heartofdepth", New)
}

func New(c core.Character, s *core.Core, count int) {
	if count >= 2 {
		var m [core.EndStatType]float64
		m[core.HydroP] = 0.15
		c.AddMod(core.CharStatMod{
			Key: "hod-2pc",
			Amount: func(a core.AttackTag) ([core.EndStatType]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		key := fmt.Sprintf("%v-hod-4pc", c.Name())

		s.Events.Subscribe(core.PostSkill, func(args ...interface{}) bool {
			s.Status.AddStatus(key, 900) //activate for 15 seoncds
			return false
		}, fmt.Sprintf("hod4-%v", c.Name()))

		c.AddMod(core.CharStatMod{
			Key: "hod-4pc",
			Amount: func(ds core.AttackTag) ([core.EndStatType]float64, bool) {
				var m [core.EndStatType]float64
				m[core.DmgP] = 0.3
				if s.Status.Duration(key) == 0 {
					return m, false
				}
				if ds != core.AttackTagNormal && ds != core.AttackTagExtra {
					return m, false
				}
				return m, true
			},
			Expiry: -1,
		})
	}
	//add flat stat to char
}
