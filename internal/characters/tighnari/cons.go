package tighnari

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c1() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = 0.15
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("tighnari-c1", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != combat.AttackTagExtra {
				return nil, false
			}
			return m, true
		},
	})
}

func (c *char) c2() {
	// crutch
	m := make([]float64, attributes.EndStatType)
	m[attributes.DendroP] = .2
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("tighnari-c2", 20*60-skillHitmark), // 12+8
		AffectedStat: attributes.DendroP,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}
