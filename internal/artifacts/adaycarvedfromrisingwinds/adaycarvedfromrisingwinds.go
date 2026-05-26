package adaycarvedfromrisingwinds

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
	core.RegisterSetFunc(keys.ADayCarvedFromRisingWinds, NewSet)
}

type Set struct {
	Index int
	Count int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

func NewSet(core *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{Count: count}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ATKP] = 0.18
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("a-day-carved-from-rising-winds-2pc", -1),
			AffectedStat: attributes.ATKP,
			Amount: func() []float64 {
				return m
			},
		})
	}

	if count >= 4 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ATKP] = 0.25

		core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
			atk, ok := args[1].(*info.AttackEvent)
			if !ok {
				return
			}
			if atk.Info.ActorIndex != char.Index() {
				return
			}
			if atk.Info.AttackTag != attacks.AttackTagNormal &&
				atk.Info.AttackTag != attacks.AttackTagExtra &&
				atk.Info.AttackTag != attacks.AttackTagElementalArt &&
				atk.Info.AttackTag != attacks.AttackTagElementalArtHold &&
				atk.Info.AttackTag != attacks.AttackTagElementalBurst {
				return
			}

			abil := "blessing-of-pastoral-winds"

			if char.IsHexerei {
				m[attributes.CR] = 0.2
				abil = "resolve-of-pastoral-winds"
			}

			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBase(abil, 6*60),
				AffectedStat: attributes.NoStat,
				Amount: func() []float64 {
					return m
				},
			})
		}, fmt.Sprintf("a-day-carved-from-rising-winds-4pc-%v", char.Base.Key.String()))
	}

	return &s, nil
}
