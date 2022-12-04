package lisa

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

func (c *char) c6() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(_ ...interface{}) bool {
		if c.Core.F < c.c6icd && c.c6icd != 0 {
			return false
		}
		if c.Core.Player.Active() == c.Index {
			//swapped to lisa

			// Create a "fake attack" to apply conductive stacks to all nearby opponents
			// Needed to ensure hitboxes are properly accounted for
			// Similar to current "Freeze Breaking" solution
			ai := combat.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Lisa C6 Conductive Status Application",
				AttackTag:  combat.AttackTagNone,
				ICDTag:     combat.ICDTagNone,
				ICDGroup:   combat.ICDGroupDefault,
				StrikeType: combat.StrikeTypeDefault,
				Element:    attributes.NoElement,
				DoNotLog:   true,
			}
			cb := func(a combat.AttackCB) {
				t, ok := a.Target.(*enemy.Enemy)
				if !ok {
					return
				}
				t.SetTag(conductiveTag, 3)
			}
			// TODO: No idea what the exact radius of this is
			//per Nosi's notes: Furthermore, the Radius of Lisa's C6 is 5m, both when in combat or not.
			c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5), -1, 0, cb)

			c.c6icd = c.Core.F + 300
		}
		return false
	}, "lisa-c6")
}
