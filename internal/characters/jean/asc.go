package jean

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

// Hits by Jean's Normal Attacks have a 50% chance to regenerate HP equal to 15% of Jean's ATK for all party members.
func (c *char) a1() combat.AttackCBFunc {
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true
		snap := a.AttackEvent.Snapshot
		if c.Core.Rand.Float64() < 0.5 {
			heal := 0.15 * (snap.BaseAtk*(1+snap.Stats[attributes.ATKP]) + snap.Stats[attributes.ATK])
			c.Core.Player.Heal(player.HealInfo{
				Caller:  c.Index,
				Target:  -1,
				Message: "Wind Companion",
				Src:     heal,
				Bonus:   c.Stat(attributes.Heal),
			})
		}
	}
}

// Using Dandelion Breeze will regenerate 20% of its Energy.
func (c *char) a4() {
	c.Energy = 16
}
