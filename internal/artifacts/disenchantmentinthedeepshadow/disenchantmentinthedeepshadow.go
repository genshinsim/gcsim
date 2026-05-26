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

// 2pc - ATK +18%.
// 4pc - When Normal Attacks hit opponents, there is a 36% chance that it will trigger Valley Rite, which will increase Normal Attack DMG by 70% of ATK.
//
//	This effect will be dispelled 0.05s after a Normal Attack deals DMG.
//	If a Normal Attack fails to trigger Valley Rite, the odds of it triggering the next time will increase by 20%.
//	This trigger can occur once every 0.2s.
func NewSet(core *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{Count: count}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ATKP] = 0.18
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("echoes-2pc", -1),
			AffectedStat: attributes.ATKP,
			Amount: func() []float64 {
				return m
			},
		})
	}

	if count >= 4 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.CR] = 0.16

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
				if e.StatusIsActive("superconduct-phys-shred") {
					return m
				}
				return nil
			},
		})
	}

	return &s, nil
}
