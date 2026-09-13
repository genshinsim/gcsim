package heartofthefurnace

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	setKey2     = "heartofthefurnace-2pc"
	setKey4     = "heartofthefurnace-4pc"
	atkKey      = setKey4 + "-atk"
	reactionKey = setKey4 + "-reaction"
)

type Set struct {
	char  *character.CharWrapper
	core  *core.Core
	Index int
	Count int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		char:  char,
		core:  c,
		Count: count,
	}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ATKP] = 0.18

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase(setKey2, -1),
			AffectedStat: attributes.ATKP,
			Amount: func() []float64 {
				return m
			},
		})
	}

	if count < 4 {
		return &s, nil
	}

	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(atkKey, -1),
		AffectedStat: attributes.ATKP,
		Amount: func() []float64 {
			m := make([]float64, attributes.EndStatType)
			if char.StatusIsActive(atkKey) {
				m[attributes.ATKP] = 0.12
			}
			return m
		},
	})

	for _, partyChar := range s.core.Player.Chars() {
		partyChar.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase(setKey4+"-reaction-dmg", -1),
			Amount: func(ai info.AttackInfo) float64 {
				if !partyChar.StatusIsActive(reactionKey) {
					return 0
				}
				if !ai.AttackTag.IsStellarReact() {
					return 0
				}
				return 0.50
			},
		})
	}

	trigger := func() {
		s.char.AddStatus(atkKey, 12*60, true)

		for _, partyChar := range s.core.Player.Chars() {
			partyChar.AddStatus(reactionKey, 12*60, true)
		}
	}

	s.core.Events.Subscribe(
		event.OnStellarConduct,
		func(args ...any) {
			ae := args[1].(*info.AttackEvent)
			if ae.Info.ActorIndex == s.char.Index() {
				trigger()
			}
		},
		fmt.Sprintf("%s-conduct-%v", setKey4, s.char.Base.Key.String()),
	)

	s.core.Events.Subscribe(
		event.OnStellarSwirl,
		func(args ...any) {
			ae := args[1].(*info.AttackEvent)
			if ae.Info.ActorIndex == s.char.Index() {
				trigger()
			}
		},
		fmt.Sprintf("%s-swirl-%v", setKey4, s.char.Base.Key.String()),
	)

	s.core.Events.Subscribe(
		event.OnEnemyDamage,
		func(args ...any) {
			ae := args[1].(*info.AttackEvent)
			if ae.Info.ActorIndex != s.char.Index() {
				return
			}
			if !ae.Info.AttackTag.IsStellarReact() {
				return
			}
			trigger()
		},
		fmt.Sprintf("%s-dmg-%v", setKey4, s.char.Base.Key.String()),
	)

	return &s, nil
}
