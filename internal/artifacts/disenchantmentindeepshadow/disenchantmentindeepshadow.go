package disenchantmentindeepshadow

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
	core.RegisterSetFunc(keys.DisenchantmentinDeepShadow, NewSet)
}

type Set struct {
	Index int
	Count int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

// 2-Piece Bonus:
// ATK +18%.

// 4-Piece Bonus:
// Increases Superconduct Reaction DMG by 80%.
// When the wielder attacks opponents affected by Superconduct,
// this attack's CRIT Rate is increased by 16%.

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		Count: count,
	}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ATKP] = 0.18

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("deep-shadow-2pc", -1),
			AffectedStat: attributes.ATKP,
			Amount: func() []float64 {
				return m
			},
		})
	}

	if count >= 4 {
		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase("deep-shadow-4pc", -1),
			Amount: func(ai info.AttackInfo) float64 {
				if ai.AttackTag == attacks.AttackTagSuperconductDamage {
					return 0.8
				}
				return 0
			},
		})

		m := make([]float64, attributes.EndStatType)
		m[attributes.CR] = 0.16

		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("deep-shadow-4pc-cr", -1),
			Amount: func(_ *info.AttackEvent, t info.Target) []float64 {
				x, ok := t.(*enemy.Enemy)
				if !ok {
					return nil
				}

				if x.StatusIsActive("superconduct-phys-shred") {
					return m
				}

				return nil
			},
		})
	}

	return &s, nil
}
