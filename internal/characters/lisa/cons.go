package lisa

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

func (c *char) c6() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		if c.Core.F < c.c6icd && c.c6icd != 0 {
			return false
		}
		if c.Core.ActiveChar == c.CharIndex() {
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
				Element:    attributes.NoElement,
				DoNotLog:   true,
			}
			cb := func(a combat.AttackCB) {
				a.Target.SetTag(conductiveTag, 3)
			}
			// TODO: No idea what the exact radius of this is
			c.Core.Combat.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), -1, 0, cb)

			c.c6icd = c.Core.F + 300
		}
		return false
	}, "lisa-c6")
}
