package noblesse

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("noblesse oblige", New)
	core.RegisterSetFunc("noblesseoblige", New)
}

func New(c core.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.DmgP] = 0.2
		c.AddMod(core.CharStatMod{
			Key: "nob-2pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return m, a == core.AttackTagElementalBurst
			},
			Expiry: -1,
		})
	}
	if count >= 4 {

		s.Events.Subscribe(core.PostBurst, func(args ...interface{}) bool {
			// s.s.Log.Debugw("\t\tNoblesse 2 pc","frame",s.F, "name", ds.CharName, "abil", ds.AbilType)
			if s.ActiveChar != c.CharIndex() {
				return false
			}

			nob, ok := s.GetCustomFlag("nob-4pc")
			//only activate if none existing
			if s.Status.Duration("nob-4pc") == 0 || (nob == c.CharIndex() && ok) {
				s.Status.AddStatus("nob-4pc", 720)
				s.SetCustomFlag("nob-4pc", c.CharIndex())
			}

			s.Log.Debugw("noblesse 4pc proc", "frame", s.F, "event", core.LogArtifactEvent, "expiry", s.Status.Duration("nob-4pc"))
			return false

		}, fmt.Sprintf("no 4pc - %v", c.Name()))

		m := make([]float64, core.EndStatType)
		m[core.ATKP] = 0.2

		s.Events.Subscribe(core.OnInitialize, func(args ...interface{}) bool {
			for _, char := range s.Chars {
				char.AddMod(core.CharStatMod{
					Key: "nob-4pc",
					Amount: func(a core.AttackTag) ([]float64, bool) {
						if s.Status.Duration("nob-4pc") > 0 {
							return m, true
						}
						return nil, false
					},
					Expiry: -1,
				})
			}
			return true
		}, "nob-4pc-init")

	}
	//add flat stat to char
}
