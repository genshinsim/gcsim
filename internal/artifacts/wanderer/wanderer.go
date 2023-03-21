package wanderer

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/artifact"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.WanderersTroupe, NewSet)
}

type Set struct {
	char  *character.CharWrapper
	count int
	Index int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }

func (s *Set) Init() error {
	if s.count < 4 {
		return nil
	}

	switch s.char.Weapon.Class {
	case weapon.WeaponClassCatalyst:
	case weapon.WeaponClassBow:
	default:
		// don't add this mod if wrong weapon class
		return nil
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.35
	s.char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("wt-4pc", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagExtra {
				return nil, false
			}
			return m, true
		},
	})

	return nil
}

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (artifact.Set, error) {
	s := Set{
		char:  char,
		count: count,
	}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.EM] = 80
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("wt-2pc", -1),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	return &s, nil
}
