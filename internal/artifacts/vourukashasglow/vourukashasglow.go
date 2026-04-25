package vourukashasglow

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.VourukashasGlow, NewSet)
}

type Set struct {
	core  *core.Core
	Index int
	Count int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		core:  c,
		Count: count,
	}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.HPP] = 0.20
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("vg-2pc", -1),
			AffectedStat: attributes.HPP,
			Amount: func() []float64 {
				return m
			},
		})
	}

	if count >= 4 {
		counter := 0
		mStack := make([]float64, attributes.EndStatType)
		mStack[attributes.DmgP] = 0.08
		addStackMod := func(idx int, duration int) {
			char.AddAttackMod(character.AttackMod{
				Base: modifier.NewBaseWithHitlag(fmt.Sprintf("vg-4pc-%v-stack", idx+1), duration),
				Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
					switch atk.Info.AttackTag {
					case attacks.AttackTagElementalArt,
						attacks.AttackTagElementalArtHold,
						attacks.AttackTagElementalBurst:
						return mStack
					default:
						return nil
					}
				},
			})
		}
		c.Events.Subscribe(event.OnPlayerHPDrain, func(args ...any) {
			di := args[0].(*info.DrainInfo)
			if di.ActorIndex != char.Index() {
				return
			}
			if di.Amount <= 0 {
				return
			}
			if !di.External {
				return
			}
			addStackMod(counter, 300)
			counter = (counter + 1) % 5
		}, fmt.Sprintf("vg-4pc-%v", char.Base.Key.String()))

		mBase := make([]float64, attributes.EndStatType)
		mBase[attributes.DmgP] = 0.1
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("vg-4pc", -1),
			Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
				switch atk.Info.AttackTag {
				case attacks.AttackTagElementalArt,
					attacks.AttackTagElementalArtHold,
					attacks.AttackTagElementalBurst:
					return mBase
				default:
					return nil
				}
			},
		})
	}

	return &s, nil
}
