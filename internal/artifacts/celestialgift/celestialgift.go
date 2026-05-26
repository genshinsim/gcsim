package celestialgift

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.CelestialGift, NewSet)
}

type Set struct {
	element attributes.Element
	Index   int
	Count   int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

func NewSet(core *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{Count: count}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ER] = 0.20
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("celestial-gift-2pc", -1),
			AffectedStat: attributes.ER,
			Amount: func() []float64 {
				return m
			},
		})
	}

	if count >= 4 {
		m := make([]float64, attributes.EndStatType)
		s.element = char.Base.Element

		core.Events.Subscribe(event.OnSkill, func(args ...any) {
			if !char.IsHexerei {
				return
			}

			abil := "lights-guidance"

			m[attributes.EleToDmgP(s.element)] = 0.20

			if core.Player.GetHexereiCount() >= 2 {
				m[attributes.EleToDmgP(s.element)] = 0.40
				abil = "mortal-hymn"
			}

			for _, c := range core.Player.Chars() {
				c.AddStatMod(character.StatMod{
					Base:         modifier.NewBaseWithHitlag(abil, 20*60),
					AffectedStat: attributes.NoStat,
					Amount: func() []float64 {
						activeCharacterElement := core.Player.Chars()[core.Player.Active()].Base.Element

						if activeCharacterElement != s.element && core.Player.GetHexereiCount() >= 2 {
							m[attributes.EleToDmgP(activeCharacterElement)] = 0.40
						}

						return m
					},
				})
			}
		}, "celestial-gift-4pc")
	}

	return &s, nil
}
