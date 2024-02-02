package braveheart

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.BraveHeart, NewSet)
}

type Set struct {
	Index int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{}

	// 2 Piece: ATK +18%
	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ATKP] = 0.18
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("braveheart-2pc", -1),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
	// 4 Piece: Increases DMG by 30% against opponents with more than 50% HP.
	if count < 4 {
		return &s, nil
	}

	if !c.Combat.DamageMode {
		return &s, nil
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.30
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("braveheart-4pc", -1),
		Amount: func(_ *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			x, ok := t.(*enemy.Enemy)
			if !ok {
				return nil, false
			}
			if x.HP()/x.MaxHP() > 0.5 {
				return m, true
			}
			return nil, false
		},
	})

	return &s, nil
}
