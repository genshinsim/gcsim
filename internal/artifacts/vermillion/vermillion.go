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
		var stacks float64
		var HPicd int

		s.Events.Subscribe(core.PostBurst, func(args ...interface{}) bool {
			if s.ActiveChar != c.CharIndex() {
				return false
			}

			nob, ok := s.GetCustomFlag("verm-4pc")
			//only activate if none existing
			if s.Status.Duration("verm-4pc") == 0 || (nob == c.CharIndex() && ok) {
				s.Status.AddStatus("verm-4pc", 16*60)
				s.SetCustomFlag("verm-4pc", c.CharIndex())
				stacks = 0
			}

			s.Log.NewEvent("verm 4pc proc", core.LogArtifactEvent, c.CharIndex(), "expiry", s.Status.Duration("verm-4pc"))
			return false

		}, fmt.Sprintf("verm 4pc - %v", c.Name()))

		m := make([]float64, core.EndStatType)
		m[core.ATKP] = 0.08

		s.Events.Subscribe(core.OnInitialize, func(args ...interface{}) bool {
			c.AddMod(core.CharStatMod{
				Key: "verm-4pc",
				Amount: func() ([]float64, bool) {
					if s.Status.Duration("verm-4pc") > 0 {

						m[core.ATKP] = 0.08 + stacks*0.1
						return m, true
					}
					stacks = 0
					return nil, false
				},
				Expiry: -1,
			})
			return true
		}, "verm-4pc-init")

		s.Events.Subscribe(core.OnCharacterHurt, func(args ...interface{}) bool {
			if s.F >= HPicd && stacks < 4 && s.Status.Duration("verm-4pc") > 0 { //grants stack if conditions are met
				stacks++
				s.Log.NewEvent("Vermillion stack gained", core.LogArtifactEvent, c.CharIndex(), "stacks", stacks)
				HPicd = s.F + 48
			}
			return false
		}, "Stack-on-hurt")

		s.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
			s.Status.DeleteStatus("verm-4pc")
			stacks = 0 //resets stacks to 0 when the character swaps
			return false
		}, "char-exit")

	}
}
