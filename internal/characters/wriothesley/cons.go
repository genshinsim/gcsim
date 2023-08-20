package wriothesley

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// When using Darkgold Wolfbite, each Prosecution Edict stack from the Passive Talent
// "There Shall Be a Reckoning for Sin" will increase said ability's DMG dealt by 40%.
// You must first unlock the Passive Talent "There Shall Be a Reckoning for Sin."
func (c *char) c2() {
	if c.Base.Ascension < 4 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("wriothesley-c2", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagElementalBurst {
				return nil, false
			}
			m[attributes.DmgP] = 0.4 * float64(c.a4Stack)
			return m, true
		},
	})
}
