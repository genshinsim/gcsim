package heartofthefurnace

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	setKey2     = "furnace-2pc"
	setKey4     = "furnace-4pc"
	atkKey      = setKey4 + "-atk"
	reactionKey = setKey4 + "-react"
)

type Set struct {
	Index int
	Count int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		Count: count,
	}

	if count < 2 {
		return &s, nil
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.18

	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(setKey2, -1),
		AffectedStat: attributes.ATKP,
		Amount: func() []float64 {
			return m
		},
	})

	if count < 4 {
		return &s, nil
	}

	m2 := make([]float64, attributes.EndStatType)
	m2[attributes.ATKP] = 0.12

	atkpBuff := character.StatMod{
		Base:         modifier.NewBaseWithHitlag(atkKey, 12*60),
		AffectedStat: attributes.ATKP,
		Amount: func() []float64 {
			return m2
		},
	}

	stellarBuff := character.ReactBonusMod{
		Base: modifier.NewBaseWithHitlag(reactionKey, 12*60),
		Amount: func(ai info.AttackInfo) float64 {
			if !ai.AttackTag.IsStellar() {
				return 0
			}
			return 0.50
		},
	}

	addBuffs := func() {
		char.AddStatMod(atkpBuff)

		for _, partyChar := range c.Player.Chars() {
			partyChar.AddReactBonusMod(stellarBuff)
		}
	}

	onReaction := func(args ...any) {
		ae := args[1].(*info.AttackEvent)
		if ae.Info.ActorIndex != char.Index() {
			return
		}
		addBuffs()
	}

	onStellarDamage := func(args ...any) {
		ae := args[1].(*info.AttackEvent)
		if ae.Info.ActorIndex != char.Index() {
			return
		}
		if !ae.Info.AttackTag.IsStellar() {
			return
		}
		addBuffs()
	}

	c.Events.Subscribe(event.OnStellarConduct, onReaction, setKey4+"-"+char.Base.Key.String())
	c.Events.Subscribe(event.OnStellarSwirl, onReaction, setKey4+"-"+char.Base.Key.String())
	c.Events.Subscribe(event.OnEnemyDamage, onStellarDamage, setKey4+"-"+char.Base.Key.String())

	return &s, nil
}
