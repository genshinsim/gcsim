package disenchantmentindeepshadow

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

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
		Base:         modifier.NewBase("disenchantment-2pc", -1),
		AffectedStat: attributes.ATKP,
		Amount: func() []float64 {
			return c2Buff
		},
	})

	if count < 4 {
		return &s, nil
	}

	char.AddReactBonusMod(character.ReactBonusMod{
		Base: modifier.NewBase("disenchantment-4pc-dmg", -1),
		Amount: func(ai info.AttackInfo) float64 {
			switch ai.AttackTag {
			case attacks.AttackTagSuperconductDamage:
				return 0.8
			case attacks.AttackTagDirectStellarConduct:
				return 0.4
			}

			return 0
		},
	})

	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = 0.16
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("disenchantment-4pc-cr", -1),
		Amount: func(_ *info.AttackEvent, t info.Target) []float64 {
			e, ok := t.(*enemy.Enemy)
			if !ok {
				return nil
			}
			if e.StatusIsActive(reactable.SuperConductShredKey) {
				return m
			}

			if e.StatusIsActive(reactable.StellarConductShredKey) {
				return m
			}
			return nil
		},
	})

	core.Events.Subscribe(event.OnLunarReactionAttack, func(args ...any) {
		r, ok := args[0].(*enemy.Enemy)
		if !ok {
			return
		}

		ae, ok := args[1].(*info.AttackEvent)
		if !ok {
			return
		}

		if r.StatusIsActive(reactable.SuperConductShredKey) || r.StatusIsActive(reactable.StellarConductShredKey) {
			ae.Snapshot.Stats[attributes.CR] += 0.16
		}
	}, "disenchantment-4pc-superconduct-"+char.Base.Key.String())

	return &s, nil
}
