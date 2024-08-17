package blizzard

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.BlizzardStrayer, NewSet)
}

type Set struct {
	Index int
	Count int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{Count: count}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.CryoP] = 0.15
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("bs-2pc", -1),
			AffectedStat: attributes.CryoP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
	if count >= 4 {
		m := make([]float64, attributes.EndStatType)
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("bs-4pc", -1),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				r, ok := t.(core.Reactable)
				if !ok {
					return nil, false
				}

				// Frozen check first so we don't mistaken coexisting cryo
				if r.AuraContains(attributes.Frozen) {
					m[attributes.CR] = 0.4
					return m, true
				}
				if r.AuraContains(attributes.Cryo) {
					m[attributes.CR] = 0.2
					return m, true
				}
				return nil, false
			},
		})
	}

	return &s, nil
}
