package yaemiko

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) a4() {
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("yaemiko-a1", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			// only trigger on elemental art damage
			if atk.Info.AttackTag != combat.AttackTagElementalArt {
				return nil, false
			}
			m[attributes.DmgP] = c.Stat(attributes.EM) * 0.0015
			return m, true
		},
	})
}
