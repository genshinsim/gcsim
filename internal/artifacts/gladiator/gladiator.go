package gladiator

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.GladiatorsFinale, NewSet)
}

type Set struct {
	char  *character.CharWrapper
	Index int
	Count int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error {
	if s.Count < 4 {
		return nil
	}

	switch s.char.Weapon.Class {
	case info.WeaponClassSpear:
	case info.WeaponClassSword:
	case info.WeaponClassClaymore:
	default:
		// don't add this mod if wrong weapon class
		return nil
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.35
	s.char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("glad-4pc", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagNormal {
				return nil, false
			}
			return m, true
		},
	})

	return nil
}

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		char:  char,
		Count: count,
	}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ATKP] = 0.18
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("glad-2pc", -1),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	return &s, nil
}
