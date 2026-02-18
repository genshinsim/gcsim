package nighttimewhispersintheechoingwoods

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.NighttimeWhispersInTheEchoingWoods, NewSet)
}

type Set struct {
	Index int
	core  *core.Core
	char  *character.CharWrapper
	lastF int
	Count int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		core:  c,
		char:  char,
		lastF: 0,
		Count: count,
	}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ATKP] = 0.18
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("nighttimewhispers-2pc", -1),
			AffectedStat: attributes.ATKP,
			Amount: func() []float64 {
				return m
			},
		})
	}

	if count >= 4 {
		c.Events.Subscribe(event.OnShielded, s.OnShielded(), fmt.Sprintf("nighttimewhispers-4pc-%v", char.Base.Key.String()))
		c.Events.Subscribe(event.OnShieldBreak, s.OnShieldBreak(), fmt.Sprintf("nighttimewhispers-4pc-%v", char.Base.Key.String()))
		c.Events.Subscribe(event.OnCharacterSwap, s.OnCharacterSwap(), fmt.Sprintf("nighttimewhispers-4pc-%v", char.Base.Key.String()))
		c.Events.Subscribe(event.OnSkill, s.OnSkill(), fmt.Sprintf("nighttimewhispers-4pc-%v", char.Base.Key.String()))
	}

	return &s, nil
}

func (s *Set) OnShielded() func(args ...any) {
	return func(args ...any) {
		shd := args[0].(shield.Shield)
		if s.core.Player.Active() != s.char.Index() {
			return
		}
		if shd.Type() == shield.Crystallize {
			s.lastF = shd.Expiry()
		}
	}
}

func (s *Set) OnShieldBreak() func(args ...any) {
	return func(args ...any) {
		shd := args[0].(shield.Shield)
		if shd.Type() != shield.Crystallize {
			return
		}
		if s.core.Player.Active() != s.char.Index() {
			return
		}
		s.lastF = s.core.F + 60
	}
}

func (s *Set) OnCharacterSwap() func(args ...any) {
	return func(args ...any) {
		prev := args[0].(int)
		active := args[1].(int)
		shd := s.core.Player.Shields.Get(shield.Crystallize)
		if shd == nil {
			return
		}
		switch s.char.Index() {
		case active:
			s.lastF = shd.Expiry()
		case prev:
			s.lastF = s.core.F + 60
		}
	}
}

func (s *Set) OnSkill() func(args ...any) {
	m := make([]float64, attributes.EndStatType)
	return func(args ...any) {
		if s.core.Player.Active() != s.char.Index() {
			return
		}
		s.char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("nighttimewhispers-4pc", 10*60),
			AffectedStat: attributes.GeoP,
			Amount: func() []float64 {
				if s.core.F <= s.lastF {
					m[attributes.GeoP] = 0.2 * 2.5
				} else {
					m[attributes.GeoP] = 0.20
				}
				return m
			},
		})
	}
}
