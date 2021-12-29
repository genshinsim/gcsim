package instructor

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("instructor", New)
}

// Implements Instructor artifact set:
// 2-Piece Bonus: Increases Elemental Mastery by 80.
// 4-Piece Bonus: Upon triggering an Elemental Reaction, increases all party members' Elemental Mastery by 120 for 8s.
func New(c core.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.EM] = 80
		c.AddMod(core.CharStatMod{
			Key: "instructor-2pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		m := make([]float64, core.EndStatType)
		m[core.EM] = 120

		add := func(args ...interface{}) bool {
			atk := args[1].(*core.AttackEvent)
			// Character must be on field to proc bonus
			if s.ActiveChar != c.CharIndex() {
				return false
			}
			// Source of elemental reaction must be the character with instructor
			if atk.Info.ActorIndex != c.CharIndex() {
				return false
			}

			// Add 120 EM to all characters except the one with instructor
			for i, char := range s.Chars {
				// Skip the one with instructor
				if i == c.CharIndex() {
					continue
				}

				char.AddMod(core.CharStatMod{
					Key: "instructor-4pc",
					Amount: func(ds core.AttackTag) ([]float64, bool) {
						return m, true
					},
					Expiry: s.F + 480,
				})
			}
			return false
		}

		for i := core.EventType(core.ReactionEventStartDelim + 1); i < core.ReactionEventEndDelim; i++ {
			s.Events.Subscribe(i, add, "4ins"+c.Name())
		}
	}
}
