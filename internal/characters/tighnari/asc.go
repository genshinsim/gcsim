package tighnari

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// After Tighnari fires a Wreath Arrow, his Elemental Mastery is increased by 50 for 4s.
//
// - checks for ascension level in aimed.go to avoid queuing this up only to fail the ascension level check
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

// For every point of Elemental Mastery Tighnari possesses, his Charged Attack and Fashioner's Tanglevine Shaft DMG are increased by 0.06%.
// The maximum DMG Bonus obtainable this way is 60%.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("tighnari-a4", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagExtra && atk.Info.AttackTag != attacks.AttackTagElementalBurst {
				return nil, false
			}

			bonus := c.Stat(attributes.EM) * 0.0006
			if bonus > 0.6 {
				bonus = 0.6
			}
			m[attributes.DmgP] = bonus
			return m, true
		},
	})
}
