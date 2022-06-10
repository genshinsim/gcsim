package wanderer

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/artifact"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterSetFunc(keys.WanderersTroupe, NewSet)
}

type Set struct {
	Index int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }
func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (artifact.Set, error) {
	s := Set{}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.EM] = 80
		char.AddStatMod("wt-2pc", -1, attributes.CryoP, func() ([]float64, bool) {
			return m, true
		})
	}
	if count >= 4 {
		switch char.Weapon.Class {
		case weapon.WeaponClassCatalyst:
		case weapon.WeaponClassBow:
		default:
			// don't add this mod if wrong weapon class
			return &s, nil
		}
		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = 0.35
		char.AddAttackMod(
			"wt-4pc",
			-1,
			func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if atk.Info.AttackTag != combat.AttackTagNormal && atk.Info.AttackTag != combat.AttackTagExtra {
					return nil, false
				}
				return m, true
			},
		)
	}

	return &s, nil
}
