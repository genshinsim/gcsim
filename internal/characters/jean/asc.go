package jean

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

// Hits by Jean's Normal Attacks have a 50% chance to regenerate HP equal to 15% of Jean's ATK for all party members.
func (c *char) makeA1CB() combat.AttackCBFunc {
	if c.Base.Ascension < 1 {
		return nil
	}
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true

		snap := a.AttackEvent.Snapshot
		if c.Core.Rand.Float64() < 0.5 {
			heal := 0.15 * (snap.BaseAtk*(1+snap.Stats[attributes.ATKP]) + snap.Stats[attributes.ATK])
			c.Core.Player.Heal(info.HealInfo{
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
	if c.Base.Ascension < 4 {
		return
	}
	c.AddEnergy("jean-a4", 16)
}
