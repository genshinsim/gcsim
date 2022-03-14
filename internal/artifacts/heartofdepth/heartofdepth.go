package heartofdepth

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterSetFunc("heart of depth", New)
	core.RegisterSetFunc("heartofdepth", New)
}

func New(c coretype.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.HydroP] = 0.15
		c.AddMod(coretype.CharStatMod{
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

		s.Subscribe(core.PostSkill, func(args ...interface{}) bool {
			if s.Player.ActiveChar != c.Index()() {
				return false
			}
			s.AddStatus(key, 900) //activate for 15 seoncds
			//add stat mod here
			c.AddPreDamageMod(coretype.PreDamageMod{
				Key: "hod-4pc",
				Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
					return m, (atk.Info.AttackTag == coretype.AttackTagNormal || atk.Info.AttackTag == coretype.AttackTagExtra)
				},
				Expiry: s.Frame + 900,
			})
			return false
		}, fmt.Sprintf("hod4-%v", c.Name()))

	}
	//add flat stat to char
}
