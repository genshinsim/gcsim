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

	if count < 2 {
		return &s, nil
	}

	c2Buff := make([]float64, attributes.EndStatType)
	c2Buff[attributes.ATKP] = 0.18

	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("a-day-carved-from-rising-winds-2pc", -1),
		AffectedStat: attributes.ATKP,
		Amount: func() []float64 {
			return c2Buff
		},
	})

	if count < 4 {
		return &s, nil
	}

	c4Buff := make([]float64, attributes.EndStatType)
	c4Buff[attributes.ATKP] = 0.25

	core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		atk, ok := args[1].(*info.AttackEvent)
		if !ok {
			return
		}
		if atk.Info.ActorIndex != char.Index() {
			return
		}

		switch atk.Info.AttackTag {
		case attacks.AttackTagNormal:
		case attacks.AttackTagExtra:
		case attacks.AttackTagElementalArt:
		case attacks.AttackTagElementalArtHold:
		case attacks.AttackTagElementalBurst:
		default:
			return
		}

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("blessing-of-pastoral-winds", 6*60),
			AffectedStat: attributes.ATKP,
			Amount: func() []float64 {
				return c4Buff
			},
		})

		if !char.IsHexerei {
			return
		}

		c4BuffHexerei := make([]float64, attributes.EndStatType)
		c4BuffHexerei[attributes.CR] = 0.2

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("resolve-of-pastoral-winds", 6*60),
			AffectedStat: attributes.CR,
			Amount: func() []float64 {
				return c4BuffHexerei
			},
		})
	}, fmt.Sprintf("a-day-carved-from-rising-winds-4pc-%v", char.Base.Key.String()))

	return &s, nil
}
