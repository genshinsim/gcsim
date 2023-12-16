package nighttimewhispersintheechoingwoods

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
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
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		core: c,
		char: char,
	}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ATKP] = 0.18
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("nwew-2pc", -1),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	if count >= 4 {
		lastF := 0
		c.Events.Subscribe(event.OnShieldBreak, func(args ...interface{}) bool {
			breaked := args[0].(shield.Shield)
			shd := c.Player.Shields.Get(shield.Crystallize)
			if shd != nil {
				return false
			}
			if c.Player.Active() != char.Index {
				return false
			}
			if breaked.Type() == shield.Crystallize {
				lastF = c.F + 60
				return false
			}
			return false
		}, fmt.Sprintf("nwew-4pc-%v", char.Base.Key.String()))

		c.Events.Subscribe(event.OnShielded, func(args ...interface{}) bool {
			shd := args[0].(shield.Shield)
			if c.Player.Active() != char.Index {
				return false
			}
			if shd.Type() == shield.Crystallize {
				lastF = shd.Expiry()
				return false
			}
			return false
		}, fmt.Sprintf("nwew-4pc-%v", char.Base.Key.String()))

		c.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
			prev := args[0].(int)
			active := args[1].(int)
			shd := c.Player.Shields.Get(shield.Crystallize)
			if shd == nil {
				return false
			}
			if active == char.Index {
				lastF = shd.Expiry()
				return false
			}
			if prev == char.Index {
				lastF = c.F + 60
				return false
			}
			return false
		}, fmt.Sprintf("nwew-4pc-%v", char.Base.Key.String()))

		c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
			if c.Player.Active() != char.Index {
				return false
			}
			m := make([]float64, attributes.EndStatType)
			char.AddAttackMod(character.AttackMod{
				Base: modifier.NewBaseWithHitlag("nwew-4pc", 10*60),
				Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
					if c.F < lastF {
						m[attributes.GeoP] = 0.2 * 2.5
					} else {
						m[attributes.GeoP] = 0.20
					}
					return m, true
				},
			})
			return false
		}, fmt.Sprintf("nwew-4pc-%v", char.Base.Key.String()))
	}

	return &s, nil
}
