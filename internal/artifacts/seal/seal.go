package seal

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/artifact"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterSetFunc(keys.EmblemOfSeveredFate, NewSet)
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
		m[attributes.ER] = 0.20
		char.AddStatMod("esr-2pc", -1, attributes.ER, func() ([]float64, bool) {
			return m, true
		})
	}
	if count >= 4 {
		m := make([]float64, attributes.EndStatType)
		er := char.Stat(attributes.ER) + 1
		amt := 0.25 * er
		if amt > 0.75 {
			amt = 0.75
		}
		m[attributes.DmgP] = amt

		char.AddAttackMod("esr-4pc", -1, func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != combat.AttackTagElementalBurst {
				return nil, false
			}

			//calc er
			er := char.Stat(attributes.ER) + 1
			amt := 0.25 * er
			if amt > 0.75 {
				amt = 0.75
			}
			m[attributes.DmgP] = amt
			return m, true
		})
	}

	return &s, nil
}
