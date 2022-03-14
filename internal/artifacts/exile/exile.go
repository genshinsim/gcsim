package exile

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterSetFunc("exile", New)
	core.RegisterSetFunc("theexile", New)
}

// 2-Piece Bonus: Energy Recharge +20%.
// 4-Piece Bonus: Using an Elemental Burst regenerates 2 Energy for all party members (excluding the wearer) every 2s for 6s. This effect cannot stack.
func New(c coretype.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.ER] = .20
		c.AddMod(coretype.CharStatMod{
			Key:    "exile-2pc",
			Expiry: -1,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
	if count >= 4 {
		// TODO: does multiple exile holders extend the duration?
		s.Subscribe(core.PreBurst, func(args ...interface{}) bool {
			if s.Player.ActiveChar != c.Index() {
				return false
			}
			if s.StatusDuration("exile") > 0 {
				return false
			}
			s.AddStatus("exile", 360)

			for _, char := range s.Chars {
				this := char
				if c.Index() == this.Index() {
					continue
				}
				// 3 ticks
				for i := 120; i <= 360; i += 120 {
					this.AddTask(func() { this.AddEnergy("exile-4pc", 2) }, "exile-energy", i)
				}
			}

			return false
		}, fmt.Sprintf("exile-4pc-%v", c.Name()))
	}
}
