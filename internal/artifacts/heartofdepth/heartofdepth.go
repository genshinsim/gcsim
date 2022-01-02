package heartofdepth

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("heart of depth", New)
	core.RegisterSetFunc("heartofdepth", New)
}

func New(c core.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.HydroP] = 0.15
		c.AddMod(core.CharStatMod{
			Key: "hod-2pc",
			Amount: func() ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		key := fmt.Sprintf("%v-hod-4pc", c.Name())

		m := make([]float64, core.EndStatType)
		m[core.DmgP] = 0.3

		s.Events.Subscribe(core.PostSkill, func(args ...interface{}) bool {
			s.Status.AddStatus(key, 900) //activate for 15 seoncds
			//add stat mod here
			c.AddPreDamageMod(core.PreDamageMod{
				Key: "hod-4pc",
				Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
					return m, (atk.Info.AttackTag == core.AttackTagNormal || atk.Info.AttackTag == core.AttackTagExtra)
				},
				Expiry: s.F + 900,
			})
			return false
		}, fmt.Sprintf("hod4-%v", c.Name()))

	}
	//add flat stat to char
}
