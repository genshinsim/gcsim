package tighnari

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) a1() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = 50
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("tighnari-a1", 4*60),
		AffectedStat: attributes.EM,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}

func (c *char) a4() {
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("tighnari-a4", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != combat.AttackTagExtra && atk.Info.AttackTag != combat.AttackTagElementalBurst {
				return nil, false
			}

			m[attributes.DmgP] = c.Stat(attributes.EM) * 0.0006
			return m, true
		},
	})
}
