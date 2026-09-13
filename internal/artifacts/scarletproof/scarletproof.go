package scarletproof

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	crKey  = "scarletproof-4pc-cr"
	dmgKey = "scarletproof-4pc-dmg"
)

type Set struct {
	Index int
	Count int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{Count: count}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ATKP] = 0.18

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("scarletproof-2pc", -1),
			AffectedStat: attributes.ATKP,
			Amount: func() []float64 {
				return m
			},
		})
	}

	if count < 4 {
		return &s, nil
	}

	cr := make([]float64, attributes.EndStatType)
	cr[attributes.CR] = 0.16

	crMod := character.StatMod{
		Base:         modifier.NewBaseWithHitlag(crKey, 10*60),
		AffectedStat: attributes.CR,
		Amount: func() []float64 {
			return cr
		},
	}

	dmgMod := character.ReactBonusMod{
		Base: modifier.NewBaseWithHitlag(dmgKey, 10*60),
		Amount: func(ai info.AttackInfo) float64 {
			if ai.AttackTag == attacks.AttackTagReactionStellarSwirl {
				return 0.40
			}
			return 0
		},
	}

	c.Events.Subscribe(event.OnStellarSwirl, func(args ...any) {
		ae, ok := args[1].(*info.AttackEvent)
		if !ok {
			return
		}

		if ae.Info.ActorIndex != char.Index() {
			return
		}

		if !char.ExtendStatus(crKey, 10*60) {
			char.AddStatMod(crMod)
		}

		if !char.ExtendStatus(dmgKey, 10*60) {
			char.AddReactBonusMod(dmgMod)
		}
	}, fmt.Sprintf("scarletproof-4pc-%v", char.Base.Key.String()))

	return &s, nil
}
