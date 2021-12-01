package huskofopulentdreams

import (
	"fmt"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("huskofopulentdreams", New)
	core.RegisterSetFunc("husk of opulent dreams", New)
}

func New(c core.Character, s *core.Core, count int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.DEFP] = 0.30
		c.AddMod(core.CharStatMod{
			Key: "husk-2pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		m := make([]float64, core.EndStatType)
		stacks := 0
		stackGainICDExpiry := 0
		// Required to check for stack loss
		lastStackGain := 0
		// Source initializes at -1
		lastSwap := -1

		// Helper function to check for stack loss
		// called after every stack gain
		checkStackLoss := func() {
			if (lastStackGain + 360) >= s.F {
				return
			}
			stacks--

			s.Log.Debugw("Husk lost stack", "frame", s.F, "event", core.LogArtifactEvent, "stacks", stacks, "last_swap", lastSwap, "last_stack_change", lastStackGain)

		}

		var gainStackOfffield func(src int) func()

		gainStackOfffield = func(src int) func() {
			return func() {
				s.Log.Debugw("Husk check for off-field stack", "frame", s.F, "event", core.LogArtifactEvent, "stacks", stacks, "last_swap", lastSwap, "last_stack_change", lastStackGain, "source", src)
				if s.ActiveChar == c.CharIndex() {
					return
				}
				// Ignore if the last source was not not the most recent swap
				if lastSwap != src {
					return
				}

				if stacks < 4 {
					stacks++
				}

				s.Log.Debugw("Husk gained off-field stack", "frame", s.F, "event", core.LogArtifactEvent, "stacks", stacks, "last_swap", lastSwap, "last_stack_change", lastStackGain)

				lastStackGain = s.F

				c.AddTask(gainStackOfffield(src), "husk-4pc-off-field-gain", 180)
				c.AddTask(checkStackLoss, "husk-4pc-stack-loss-check", 360)
			}
		}

		// Initiate off-field stacking if off-field at start of the sim
		s.Events.Subscribe(core.OnInitialize, func(args ...interface{}) bool {
			if s.ActiveChar != c.CharIndex() {
				c.AddTask(gainStackOfffield(s.F), "husk-4pc-off-field-gain", 1)
			}
			return true
		}, "husk-4pc-off-field-stack-init")

		s.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
			prev := args[0].(int)
			if prev != c.CharIndex() {
				return false
			}
			lastSwap = s.F
			c.AddTask(gainStackOfffield(s.F), "husk-4pc-off-field-gain", 180)
			return false
		}, "husk-4pc-off-field-gain")

		s.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
			ds := args[1].(*core.Snapshot)
			// Only triggers when onfield
			if s.ActiveChar != c.CharIndex() {
				return false
			}
			if ds.ActorIndex != c.CharIndex() {
				return false
			}
			if stackGainICDExpiry > s.F {
				return false
			}
			if ds.Element != core.Geo {
				return false
			}

			if stacks < 4 {
				stacks++
			}

			s.Log.Debugw("Husk gained on-field stack", "frame", s.F, "event", core.LogArtifactEvent, "stacks", stacks, "last_swap", lastSwap, "last_stack_change", lastStackGain)

			lastStackGain = s.F
			stackGainICDExpiry = s.F + 18
			c.AddTask(checkStackLoss, "husk-4pc-stack-loss-check", 360)

			return false
		}, fmt.Sprintf("husk-4pc-%v", c.Name()))

		c.AddMod(core.CharStatMod{
			Key: "husk-4pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				m[core.DEFP] = 0.06 * float64(stacks)
				m[core.GeoP] = 0.06 * float64(stacks)

				return m, true
			},
			Expiry: -1,
		})

	}
}
