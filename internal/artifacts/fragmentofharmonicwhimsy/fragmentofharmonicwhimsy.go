package fragmentofharmonicwhimsy

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.FragmentOfHarmonicWhimsy, NewSet)
}

type Set struct {
	stacks int
	Index  int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }

const (
	fohw2pc = "fragmentofharmonicwhimsy-2pc"
	fohw4pc = "fragmentofharmonicwhimsy-4pc"
)

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	var s Set

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ATKP] = 0.18
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase(fohw2pc, -1),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	if count >= 4 {
		m := make([]float64, attributes.EndStatType)
		c.Events.Subscribe(event.OnHPDebt, func(args ...interface{}) bool {
			index := args[0].(int)
			amount := args[1].(float64)
			if char.Index != index || amount == 0 {
				return false
			}

			if !char.StatusIsActive(fohw4pc) {
				s.stacks = 0
			}

			if s.stacks < 3 {
				s.stacks++
			}

			m[attributes.DmgP] = 0.18 * float64(s.stacks)

			char.AddStatMod(character.StatMod{
				Base: modifier.NewBaseWithHitlag(fohw4pc, 6*60),
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})

			return false
		}, fmt.Sprintf("fragmentofharmonicwhimsy-hp-debt-%v", char.Base.Key.String()))
	}

	return &s, nil
}
