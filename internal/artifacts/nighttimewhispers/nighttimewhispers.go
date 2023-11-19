package nighttimewhispers

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

const (
	buffKey     = "nighttime-whispers-buff"
	shieldMulti = 2.5
	buffVal     = 0.16
)

func init() {
	core.RegisterSetFunc(keys.NighttimeWhispersInTheEchoingWoods, NewSet)
}

type Set struct {
	char  *character.CharWrapper
	count int
	Index int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }

func (s *Set) Init() error {
	return nil
}

func NewSet(c *core.Core, char *character.CharWrapper, count int, _ map[string]int) (info.Set, error) {
	s := Set{
		char:  char,
		count: count,
	}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ATKP] = 0.18
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("nighttime-whispers-2pc", -1),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	if count >= 4 {
		f := func(args ...interface{}) bool {
			if c.Player.Active() != char.Index {
				return false
			}
			bonus := 1.0
			// TODO: Need to make the Crystallise bonus not tied to the base buff
			m := make([]float64, attributes.EndStatType)
			char.AddStatMod(character.StatMod{
				Base: modifier.NewBaseWithHitlag(buffKey, 60*10),
				Amount: func() ([]float64, bool) {
					if char.Index == c.Player.Active() && c.Player.Shields.PlayerIsShielded() {
						s := c.Player.Shields.List()
						for _, t := range s {
							if t.Type() == shield.Crystallize {
								bonus = shieldMulti
								break
							}
						}
					}
					m[attributes.GeoP] = bonus * buffVal
					return m, true
				},
			})
			return false
		}
		c.Events.Subscribe(event.OnSkill, f, fmt.Sprintf("nighttime-whispers-4pc-%v", char.Base.Key.String()))
	}

	return &s, nil
}
