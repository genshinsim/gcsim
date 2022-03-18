package vermillion

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("vermillion hereafter", New)
	core.RegisterSetFunc("vermillionhereafter", New)
	core.RegisterSetFunc("vermillion", New)
	core.RegisterSetFunc("verm", New)
}

func New(c core.Character, s *core.Core, count int, params map[string]int) {

	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.ATKP] = 0.18
		c.AddMod(core.CharStatMod{
			Key: "verm-2pc",
			Amount: func() ([]float64, bool) {
				return m, true
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

			nob, ok := s.GetCustomFlag("verm-4pc")
			//only activate if none existing
			if s.Status.Duration("verm-4pc") == 0 || (nob == c.CharIndex() && ok) {
				s.Status.AddStatus("verm-4pc", 16*60)
				s.SetCustomFlag("verm-4pc", c.CharIndex())
			}

			s.Log.NewEvent("verm 4pc proc", core.LogArtifactEvent, c.CharIndex(), "expiry", s.Status.Duration("verm-4pc"))
			return false

		}, fmt.Sprintf("verm 4pc - %v", c.Name()))

		m := make([]float64, core.EndStatType)
		m[core.ATKP] = 0.2

		s.Events.Subscribe(core.OnInitialize, func(args ...interface{}) bool {
			c.AddMod(core.CharStatMod{
				Key: "verm-4pc",
				Amount: func() ([]float64, bool) {
					if s.Status.Duration("verm-4pc") > 0 {
						return m, true
					}
					return nil, false
				},
				Expiry: -1,
			})
			return true
		}, "verm-4pc-init")

	}
	//add flat stat to char
}
