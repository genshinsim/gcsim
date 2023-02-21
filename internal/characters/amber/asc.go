package amber

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// Increases the CRIT Rate of Fiery Rain by 10% and widens its AoE by 30%.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	// crit
	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = .1
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("amber-a1", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			return m, atk.Info.AttackTag == combat.AttackTagElementalBurst
		},
	})
	// AoE
	c.burstRadius *= 1.3
}

// Aimed Shot hits on weak points increase ATK by 15% for 10s.
func (c *char) makeA4CB() combat.AttackCBFunc {
	if c.Base.Ascension < 4 {
		return nil
	}
	done := false
	return func(a combat.AttackCB) {
		if !a.AttackEvent.Info.HitWeakPoint {
			return
		}
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true

		m := make([]float64, attributes.EndStatType)
		m[attributes.ATKP] = 0.15
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("amber-a4", 600),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
}
