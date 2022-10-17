package candace

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const a4Key = "candace-a4"

func (c *char) a4(char *character.CharWrapper) {
	m := make([]float64, attributes.EndStatType)
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(a4Key, -1),
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			if !c.StatusIsActive(burstKey) {
				return nil, false
			}
			if atk.Info.AttackTag != combat.AttackTagNormal {
				return nil, false
			}
			if atk.Info.Element == attributes.Physical || atk.Info.Element == attributes.NoElement {
				return nil, false
			}
			m[attributes.DmgP] = 0.005 * c.MaxHP() / 1000
			return m, true
		},
	})
}
