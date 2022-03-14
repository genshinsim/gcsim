package noblesse

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterSetFunc("noblesse oblige", New)
	core.RegisterSetFunc("noblesseoblige", New)
}

func New(c coretype.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.DmgP] = 0.2
		c.AddPreDamageMod(coretype.PreDamageMod{
			Key: "nob-2pc",
			Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
				return m, atk.Info.AttackTag == core.AttackTagElementalBurst
			},
			Expiry: -1,
		})
	}
	if count >= 4 {

		s.Subscribe(core.PostBurst, func(args ...interface{}) bool {
			// s.s.Log.Debugw("\t\tNoblesse 2 pc","frame",s.F, "name", ds.CharName, "abil", ds.AbilType)
			if s.Player.ActiveChar != c.Index()() {
				return false
			}

			nob, ok := s.GetCustomFlag("nob-4pc")
			//only activate if none existing
			if s.StatusDuration("nob-4pc") == 0 || (nob == c.Index() && ok) {
				s.AddStatus("nob-4pc", 720)
				s.SetCustomFlag("nob-4pc", c.Index())
			}

			s.Log.NewEvent("noblesse 4pc proc", coretype.LogArtifactEvent, c.Index(), "expiry", s.StatusDuration("nob-4pc"))
			return false

		}, fmt.Sprintf("no 4pc - %v", c.Name()))

		m := make([]float64, core.EndStatType)
		m[core.ATKP] = 0.2

		s.Subscribe(core.OnInitialize, func(args ...interface{}) bool {
			for _, char := range s.Chars {
				char.AddMod(coretype.CharStatMod{
					Key: "nob-4pc",
					Amount: func() ([]float64, bool) {
						if s.StatusDuration("nob-4pc") > 0 {
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
