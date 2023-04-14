package vourukashasglow

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/artifact"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.VourukashasGlow, NewSet)
}

type Set struct {
	icd   int
	core  *core.Core
	Index int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }

func (s *Set) Init() error { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (artifact.Set, error) {
	s := Set{
		core: c,
	}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.HPP] = 0.20
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("dew-2pc", -1),
			AffectedStat: attributes.HPP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	if count >= 4 {
		counter := 0
		permStacks := param["stacks"]
		if permStacks > 5 {
			permStacks = 5
		}
		mStack := make([]float64, attributes.EndStatType)
		mStack[attributes.DmgP] = 0.08
		addStackMod := func(idx int, duration int) {
			char.AddAttackMod(character.AttackMod{
				Base: modifier.NewBaseWithHitlag(fmt.Sprintf("dew-4pc-%v-stack", idx+1), duration),
				Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
					switch atk.Info.AttackTag {
					case attacks.AttackTagElementalArt,
						attacks.AttackTagElementalArtHold,
						attacks.AttackTagElementalBurst:
						return mStack, true
					default:
						return nil, false
					}
				},
			})
		}
		for i := 0; i < permStacks; i++ {
			addStackMod(i, -1)
		}
		c.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
			di := args[0].(player.DrainInfo)
			if di.Amount <= 0 {
				return false
			}
			if di.ActorIndex != char.Index {
				return false
			}
			if counter >= permStacks {
				addStackMod(counter, 300)
			}
			counter = (counter + 1) % 5
			return false
		}, fmt.Sprintf("dew-4pc-%v", char.Base.Key.String()))

		mBase := make([]float64, attributes.EndStatType)
		mBase[attributes.DmgP] = 0.1
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("dew-4pc", -1),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				switch atk.Info.AttackTag {
				case attacks.AttackTagElementalArt,
					attacks.AttackTagElementalArtHold,
					attacks.AttackTagElementalBurst:
					return mBase, true
				default:
					return nil, false
				}
			},
		})
	}

	return &s, nil
}
