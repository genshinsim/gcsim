package albedo

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// c4: active member +30% plunge attack in skill field
func (c *char) c4() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.3
	for _, char := range c.Core.Player.Chars() {
		this := char
		char.AddAttackMod(character.AttackMod{Base: modifier.NewBase("albedo-c4", -1), Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if c.Core.Player.Active() != this.Index {
				return nil, false
			}
			if atk.Info.AttackTag != combat.AttackTagPlunge {
				return nil, false
			}
			if c.Tags["elevator"] != 1 {
				return nil, false
			}
			return m, true
		}})
	}
}

// c6: active protected by crystallize +17% dmg
func (c *char) c6() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.17
	c.AddStatMod(character.StatMod{Base: modifier.NewBase("albedo-c6", -1), AffectedStat: attributes.DmgP, Amount: func() ([]float64, bool) {
		if c.Tags["elevator"] != 1 {
			return nil, false
		}
		if c.Core.Player.Shields.Get(shield.ShieldCrystallize) == nil {
			return nil, false
		}
		return m, true
	}})
}
