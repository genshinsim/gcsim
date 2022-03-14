package instructor

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterSetFunc("instructor", New)
}

// Implements Instructor artifact set:
// 2-Piece Bonus: Increases Elemental Mastery by 80.
// 4-Piece Bonus: Upon triggering an Elemental Reaction, increases all party members' Elemental Mastery by 120 for 8s.
func New(c coretype.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.EM] = 80
		c.AddMod(coretype.CharStatMod{
			Key: "instructor-2pc",
			Amount: func() ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		m := make([]float64, core.EndStatType)
		m[core.EM] = 120

		add := func(args ...interface{}) bool {
			atk := args[1].(*coretype.AttackEvent)
			// Character must be on field to proc bonus
			if s.Player.ActiveChar != c.Index()() {
				return false
			}
			// Source of elemental reaction must be the character with instructor
			if atk.Info.ActorIndex != c.Index() {
				return false
			}

			// Add 120 EM to all characters except the one with instructor
			for i, char := range s.Chars {
				// Skip the one with instructor
				if i == c.Index() {
					continue
				}

				char.AddMod(coretype.CharStatMod{
					Key: "instructor-4pc",
					Amount: func() ([]float64, bool) {
						return m, true
					},
					Expiry: s.Frame + 480,
				})
			}
			return false
		}

		for i := core.EventType(core.ReactionEventStartDelim + 1); i < core.ReactionEventEndDelim; i++ {
			s.Subscribe(i, add, "4ins"+c.Name())
		}
	}
}
