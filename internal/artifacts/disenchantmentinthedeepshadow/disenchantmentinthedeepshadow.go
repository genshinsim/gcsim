package disenchantmentinthedeepshadow

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

func init() {
	core.RegisterSetFunc(keys.DisenchantmentInTheDeepShadow, NewSet)
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
		Base:         modifier.NewBase("disenchantment-2pc", -1),
		AffectedStat: attributes.ATKP,
		Amount: func() []float64 {
			return c2Buff
		},
	})

	if count < 4 {
		return &s, nil
	}

	c4Buff := make([]float64, attributes.EndStatType)
	c4Buff[attributes.CR] = 0.16

	char.AddReactBonusMod(character.ReactBonusMod{
		Base: modifier.NewBase("disenchantment-reaction-bonus-4pc", -1),
		Amount: func(ai info.AttackInfo) float64 {
			if ai.AttackTag != attacks.AttackTagSuperconductDamage {
				return 0
			}

			return 0.8
		},
	})

	// TODO: check if it works while character is off-field

	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("disenchantment-on-superconduct-4pc", -1),
		Amount: func(_ *info.AttackEvent, t info.Target) []float64 {
			e, ok := t.(*enemy.Enemy)
			if !ok {
				return nil
			}
			if e.StatusIsActive(reactable.SuperConductShredKey) {
				return c4Buff
			}
			return nil
		},
	})

	return &s, nil
}
