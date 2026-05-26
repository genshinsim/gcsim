package celestialgift

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

const (
	duration = 60 * 20
	set4Key  = "celestialgift-4pc"
)

func init() {
	core.RegisterSetFunc(keys.CelestialGift, NewSet)
}

type Set struct {
	Index int
	Count int
	buff  []float64
	core  *core.Core
	char  *character.CharWrapper
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error {
	c := s.core
	char := s.char
	if s.Count < 4 || !char.IsHexerei {
		return nil
	}

	buffActive := func() bool {
		for _, c := range c.Player.Chars() {
			if c.StatusIsActive(set4Key) {
				return true
			}
		}
		return false
	}

	buffActiveOfElement := func(e attributes.Element) bool {
		for _, c := range c.Player.Chars() {
			if c.StatusIsActive(set4Key) && c.Base.Element == e {
				return true
			}
		}
		return false
	}

	elems := map[attributes.Element]bool{}
	for _, x := range c.Player.Chars() {
		elems[x.Base.Element] = true
	}

	buffStrength := 0.2
	isHexereiSecretRite := c.Player.GetHexereiCount() >= 2
	if isHexereiSecretRite {
		buffStrength = 0.4
	}
	for elem := range elems {
		buffedDmgP := attributes.EleToDmgP(elem)
		for _, x := range c.Player.Chars() {
			x.AddStatMod(character.StatMod{
				Base:         modifier.NewBase(set4Key+"-"+elem.String(), -1),
				AffectedStat: buffedDmgP,
				Amount: func() []float64 {
					clear(s.buff)
					s.buff[buffedDmgP] = buffStrength
					if buffActiveOfElement(elem) {
						return s.buff
					}
					if isHexereiSecretRite && c.Player.ActiveChar().Base.Element == elem && buffActive() {
						return s.buff
					}
					return nil
				},
			})
		}
	}

	return nil
}

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		Count: count,
		buff:  make([]float64, attributes.EndStatType),
		core:  c,
		char:  char,
	}

	if count < 2 {
		return &s, nil
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.ER] = 0.2
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("celestialgift-2pc", -1),
		AffectedStat: attributes.ER,
		Amount: func() []float64 {
			return m
		},
	})

	if count < 4 || !char.IsHexerei {
		return &s, nil
	}

	// NOTE: Description says "The equipping character can trigger this effect while off-field."
	// but currently no characters can use skills while off field, so this is ignored
	c.Events.Subscribe(event.OnSkill, func(args ...any) {
		if c.Player.Active() != char.Index() {
			return
		}
		char.AddStatus(set4Key, duration, true)
	}, fmt.Sprintf(set4Key+"-%v", char.Base.Key.String()))

	return &s, nil
}
