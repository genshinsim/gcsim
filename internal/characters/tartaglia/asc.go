package tartaglia

import (
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

// Extends Riptide duration by 8s.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	c.riptideDuration += 8 * 60
}

// When Tartaglia is in Foul Legacy: Raging Tide's Melee Stance, on dealing a CRIT hit,
// Normal and Charged Attacks apply the Riptide status effect to opponents.
func (c *char) makeA4CB() combat.AttackCBFunc {
	if c.Base.Ascension < 4 {
		return nil
	}
	return func(a combat.AttackCB) {
		if a.IsCrit {
			t, ok := a.Target.(*enemy.Enemy)
			if !ok {
				return
			}
			c.applyRiptide("melee", t)
		}
	}
}
