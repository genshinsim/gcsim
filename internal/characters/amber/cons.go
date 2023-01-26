package amber

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// C2
// Explosion via manual detonation deals 200% additional DMG.
func (c *char) c2() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 2
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag("amber-c2", -1),
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != combat.AttackTagElementalArt && atk.Info.Abil != manualExplosionAbil {
				return nil, false
			}
			return m, true
		},
	})
}
