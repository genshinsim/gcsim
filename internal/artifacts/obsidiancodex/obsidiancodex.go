package obsidiancodex

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/template/nightsoul"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.ObsidianCodex, NewSet)
}

type Set struct {
	Index        int
	Count        int
	consumeCount float64
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{Count: count}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = 0.15
		char.AddStatMod(character.StatMod{
			Base: modifier.NewBase("obsidiancodex-2pc", -1),
			Amount: func() ([]float64, bool) {
				if !char.StatusIsActive(nightsoul.NightsoulBlessingStatus) {
					return nil, false
				}
				if c.Player.Active() != char.Index {
					return nil, false
				}
				return m, true
			},
		})
	}

	if count >= 4 {
		const icdKey = "obsidiancodex-4pc-icd"
		m := make([]float64, attributes.EndStatType)
		m[attributes.CR] = 0.4
		c.Events.Subscribe(event.OnNightsoulConsume, func(args ...interface{}) bool {
			idx := args[0].(int)
			amount := args[1].(float64)
			if char.Index != idx {
				return false
			}
			if c.Player.Active() != char.Index {
				return false
			}
			if char.StatusIsActive(icdKey) {
				return false
			}
			char.AddStatus(icdKey, 60, true)
			s.consumeCount += amount
			if s.consumeCount >= 1 {
				s.consumeCount = 0
				char.AddStatMod(character.StatMod{
					Base: modifier.NewBaseWithHitlag("obsidiancodex-4pc", 6*60),
					Amount: func() ([]float64, bool) {
						return m, true
					},
				})
			}
			return false
		}, fmt.Sprintf("obsidiancodex-4pc-%v", char.Base.Key.String()))
	}

	return &s, nil
}
