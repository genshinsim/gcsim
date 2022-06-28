package yelan

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) a1() {
	partyEleTypes := make(map[attributes.Element]bool)
	for _, char := range c.Core.Player.Chars() {
		partyEleTypes[char.Base.Element] = true
	}
	count := len(partyEleTypes)

	m := make([]float64, attributes.EndStatType)
	m[attributes.HPP] = float64(count) * 0.06
	if count >= 4 {
		m[attributes.HPP] = 0.3
	}

	c.AddStatMod(character.StatMod{Base: modifier.NewBase("yelan-a1", -1), AffectedStat: attributes.HPP, Amount: func() ([]float64, bool) {
		return m, true
	}})
}

func (c *char) a4() {
	started := c.Core.F
	m := make([]float64, attributes.EndStatType)
	for _, char := range c.Core.Player.Chars() {
		this := char
		this.AddAttackMod(character.AttackMod{Base: modifier.NewBase("yelan-a4", 15*60), Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			//char must be active
			if c.Core.Player.Active() != this.Index {
				return nil, false
			}
			//floor time elapsed
			dmg := float64(int((c.Core.F-started)/60))*0.035 + 0.01
			if dmg > 0.5 {
				dmg = 0.5
			}
			m[attributes.DmgP] = dmg
			return m, true
		}})
	}
}
