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
			Base:         modifier.NewBase("nwitew-2pc", -1),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	if count >= 4 {
		// TODO need better approach
		lastF := 0
		c.Events.Subscribe(event.OnTick, func(args ...interface{}) bool {
			shd := c.Player.Shields.Get(shield.Crystallize)
			if c.Player.Active() != char.Index {
				return false
			}
			if shd != nil {
				lastF = c.F
			}
			return false
		}, "")

		c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
			if c.Player.Active() != char.Index {
				return false
			}
			m := make([]float64, attributes.EndStatType)
			m[attributes.GeoP] = 0.20
			if c.F < lastF+60 {
				m[attributes.GeoP] = 0.2 * 2.5
			}
			char.AddStatMod(character.StatMod{
				Base: modifier.NewBaseWithHitlag("nwitew-4pc", 10*60),
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})
			return false
		}, fmt.Sprintf("nwitew-4pc-skill-%v", char.Base.Key.String()))
	}

	return &s, nil
}
