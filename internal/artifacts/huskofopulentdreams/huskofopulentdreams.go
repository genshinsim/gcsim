package huskofopulentdreams

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterSetFunc("huskofopulentdreams", New)
	core.RegisterSetFunc("husk of opulent dreams", New)
}

/**
A character equipped with this Artifact set will obtain the Curiosity effect in the following conditions:
When on the field, the character gains 1 stack after hitting an opponent with a Geo attack,
triggering a maximum of once every 0.3s. When off the field, the character gains 1 stack every 3s.

Curiosity can stack up to 4 times, each providing 6% DEF and a 6% Geo DMG Bonus. When 6 seconds pass
without gaining a Curiosity stack, 1 stack is lost.
**/
func New(c coretype.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.DEFP] = 0.30
		c.AddMod(coretype.CharStatMod{
			Key: "husk-2pc",
			Amount: func() ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		m := make([]float64, core.EndStatType)
		stacks := params["stacks"]
		stackGainICDExpiry := 0
		// Required to check for stack loss
		lastStackGain := 0
		// Source initializes at -1
		lastSwap := -1

		// Helper function to check for stack loss
		// called after every stack gain
		checkStackLoss := func() {
			if (lastStackGain + 360) >= s.Frame {
				return
			}
			stacks--
			s.Log.NewEvent("Husk lost stack", coretype.LogArtifactEvent, c.Index(), "stacks", stacks, "last_swap", lastSwap, "last_stack_change", lastStackGain)

		}

		var gainStackOfffield func(src int) func()

		gainStackOfffield = func(src int) func() {
			return func() {
				s.Log.NewEvent("Husk check for off-field stack", coretype.LogArtifactEvent, c.Index(), "stacks", stacks, "last_swap", lastSwap, "last_stack_change", lastStackGain, "source", src)
				if s.ActiveChar == c.Index() {
					return
				}
				// Ignore if the last source was not not the most recent swap
				if lastSwap != src {
					return
				}

				if stacks < 4 {
					stacks++
				}

				s.Log.NewEvent("Husk gained off-field stack", coretype.LogArtifactEvent, c.Index(), "stacks", stacks, "last_swap", lastSwap, "last_stack_change", lastStackGain)

				lastStackGain = s.Frame

				c.AddTask(gainStackOfffield(src), "husk-4pc-off-field-gain", 180)
				c.AddTask(checkStackLoss, "husk-4pc-stack-loss-check", 360)
			}
		}

		// Initiate off-field stacking if off-field at start of the sim
		s.Subscribe(core.OnInitialize, func(args ...interface{}) bool {
			if s.Player.ActiveChar != c.Index()() {
				c.AddTask(gainStackOfffield(s.Frame), "husk-4pc-off-field-gain", 1)
			}
			return true
		}, "husk-4pc-off-field-stack-init")

		s.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
			prev := args[0].(int)
			if prev != c.Index() {
				return false
			}
			lastSwap = s.Frame
			c.AddTask(gainStackOfffield(s.Frame), "husk-4pc-off-field-gain", 180)
			return false
		}, "husk-4pc-off-field-gain")

		s.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
			atk := args[1].(*coretype.AttackEvent)
			// Only triggers when onfield
			if s.Player.ActiveChar != c.Index()() {
				return false
			}
			if atk.Info.ActorIndex != c.Index() {
				return false
			}
			if stackGainICDExpiry > s.Frame {
				return false
			}
			if atk.Info.Element != core.Geo {
				return false
			}

			if stacks < 4 {
				stacks++
			}

			s.Log.NewEvent("Husk gained on-field stack", coretype.LogArtifactEvent, c.Index(), "stacks", stacks, "last_swap", lastSwap, "last_stack_change", lastStackGain)

			lastStackGain = s.Frame
			stackGainICDExpiry = s.Frame + 18
			c.AddTask(checkStackLoss, "husk-4pc-stack-loss-check", 360)

			return false
		}, fmt.Sprintf("husk-4pc-%v", c.Name()))

		c.AddMod(coretype.CharStatMod{
			Key: "husk-4pc",
			Amount: func() ([]float64, bool) {
				m[core.DEFP] = 0.06 * float64(stacks)
				m[core.GeoP] = 0.06 * float64(stacks)

				return m, true
			},
			Expiry: -1,
		})

	}
}
