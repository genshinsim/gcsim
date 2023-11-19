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
	buffKey       = "nighttime-whispers-buff"
	buffVal       = 0.16
	secondBuffKey = "nighttime-whispers-second-buff"
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

	// 2pc - ATK +18%.
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

	// 4pc - After using an Elemental Skill, gain a 16% Geo DMG Bonus for 10s
	if count >= 4 {
		f := func(args ...interface{}) bool {
			if c.Player.Active() != char.Index {
				return false
			}
			m := make([]float64, attributes.EndStatType)
			char.AddStatMod(character.StatMod{
				Base: modifier.NewBaseWithHitlag(buffKey, 60*10),
				Amount: func() ([]float64, bool) {
					m[attributes.GeoP] = buffVal
					return m, true
				},
			})

			// While under a shield granted by the Crystallize reaction, the above
			// effect will be increased by 150%, and this additional increase disappears
			// 1s after that shield is lost.
			for i := 0; i < 11*60; i += 30 { // An extra second to account for possible hitlag extension
				char.QueueCharTask(
					func() {
						// Checks that base buff is active
						if !char.StatusIsActive(buffKey) {
							if char.StatusIsActive(secondBuffKey) {
								char.RemoveTag(secondBuffKey)
							}
							return
						}

						// Checks for a Crystallise Shield.
						if char.Index == c.Player.Active() && c.Player.Shields.PlayerIsShielded() {
							s := c.Player.Shields.List()
							for _, t := range s {
								if t.Type() == shield.Crystallize {
									n := make([]float64, attributes.EndStatType)
									char.AddStatMod(character.StatMod{
										Base: modifier.NewBaseWithHitlag(secondBuffKey, 60),
										Amount: func() ([]float64, bool) {
											n[attributes.GeoP] = 1.5 * buffVal
											return n, true
										},
									})
									break
								}
							}
						}
						return
					},
					i,
				)
			}
			return false
		}

		c.Events.Subscribe(event.OnSkill, f, fmt.Sprintf("nighttime-whispers-4pc-%v", char.Base.Key.String()))
	}

	return &s, nil
}
