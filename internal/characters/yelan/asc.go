package yelan

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// When the party has 1/2/3/4 Elemental Types, Yelan's Max HP is increased by 6%/12%/18%/30%.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
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

	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("yelan-a1", -1),
		AffectedStat: attributes.HPP,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}

// So long as an Exquisite Throw is in play, your own active character deals 1% more DMG.
// This increases by a further 3.5% DMG every second. The maximum increase to DMG dealt is 50%.
// The pre-existing effect will be dispelled if Depth-Clarion Dice is recast during its duration.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	started := c.Core.F
	for _, char := range c.Core.Player.Chars() {
		this := char
		this.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("yelan-a4", 15*60),
			Amount: func(_ *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
				//char must be active
				if c.Core.Player.Active() != this.Index {
					return nil, false
				}
				//floor time elapsed
				dmg := float64(int((c.Core.F-started)/60))*0.035 + 0.01
				if dmg > 0.5 {
					dmg = 0.5
				}
				c.a4buff[attributes.DmgP] = dmg
				return c.a4buff, true
			},
		})
	}
}
