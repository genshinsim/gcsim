package instructor

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.Instructor, NewSet)
}

type Set struct {
	Index int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }

// Implements Instructor artifact set:
// 2-Piece Bonus: Increases Elemental Mastery by 80.
// 4-Piece Bonus: Upon triggering an Elemental Reaction, increases all party members' Elemental Mastery by 120 for 8s.
func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.EM] = 80
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("instructor-2pc", -1),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
	if count >= 4 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.EM] = 120

		// TODO: does multiple instructor holders extend the duration?
		add := func(args ...interface{}) bool {
			atk := args[1].(*combat.AttackEvent)
			// Character must be on field to proc bonus
			if c.Player.Active() != char.Index {
				return false
			}
			// Source of elemental reaction must be the character with instructor
			if atk.Info.ActorIndex != char.Index {
				return false
			}

			// Add 120 EM to all characters
			for _, this := range c.Player.Chars() {
				this.AddStatMod(character.StatMod{
					Base:         modifier.NewBaseWithHitlag("instructor-4pc", 480),
					AffectedStat: attributes.EM,
					Amount: func() ([]float64, bool) {
						return m, true
					},
				})
			}
			return false
		}

		for i := event.ReactionEventStartDelim + 1; i < event.OnShatter; i++ {
			c.Events.Subscribe(i, add, fmt.Sprintf("instructor-4pc-%v", char.Base.Key.String()))
		}
	}

	return &s, nil
}
