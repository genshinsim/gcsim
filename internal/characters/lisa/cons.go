package lisa

import (
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

func (c *char) c6() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(_ ...interface{}) bool {
		if c.Core.F < c.c6icd && c.c6icd != 0 {
			return false
		}
		if c.Core.Player.Active() == c.Index {
			//swapped to lisa
			enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6), nil)
			for _, e := range enemies {
				e.SetTag(conductiveTag, 3)
			}
			c.c6icd = c.Core.F + 300
		}
		return false
	}, "lisa-c6")
}
