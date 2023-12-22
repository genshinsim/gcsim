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
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		core:  c,
		char:  char,
		lastF: 0,
	}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ATKP] = 0.18
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("nighttimewhispers-2pc", -1),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return m, true
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

func (s *Set) OnShielded() func(args ...interface{}) bool {
	return func(args ...interface{}) bool {
		shd := args[0].(shield.Shield)
		if s.core.Player.Active() != s.char.Index {
			return false
		}
		if shd.Type() == shield.Crystallize {
			s.lastF = shd.Expiry()
			return false
		}
		return false
	}
}

func (s *Set) OnShieldBreak() func(args ...interface{}) bool {
	return func(args ...interface{}) bool {
		shd := args[0].(shield.Shield)
		if shd.Type() != shield.Crystallize {
			return false
		}
		if s.core.Player.Active() != s.char.Index {
			return false
		}
		s.lastF = s.core.F + 60
		return false
	}
}

func (s *Set) OnCharacterSwap() func(args ...interface{}) bool {
	return func(args ...interface{}) bool {
		prev := args[0].(int)
		active := args[1].(int)
		shd := s.core.Player.Shields.Get(shield.Crystallize)
		if shd == nil {
			return false
		}
		if active == s.char.Index {
			s.lastF = shd.Expiry()
			return false
		}
		if prev == s.char.Index {
			s.lastF = s.core.F + 60
			return false
		}
		return false
	}
}

func (s *Set) OnSkill() func(args ...interface{}) bool {
	return func(args ...interface{}) bool {
		if s.core.Player.Active() != s.char.Index {
			return false
		}
		m := make([]float64, attributes.EndStatType)
		s.char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("nighttimewhispers-4pc", 10*60),
			AffectedStat: attributes.GeoP,
			Amount: func() ([]float64, bool) {
				if s.core.F <= s.lastF {
					m[attributes.GeoP] = 0.2 * 2.5
				} else {
					m[attributes.GeoP] = 0.20
				}
				return m, true
			},
		})
		return false
	}
}
