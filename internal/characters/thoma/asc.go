package thoma

import (
	"github.com/genshinsim/gcsim/pkg/core/event"
)

func (c *char) a1() {
	c.Core.Player.Shields.AddShieldBonusMod("thoma-a1", -1, func() (float64, bool) {
		if c.Tags["shielded"] == 0 {
			return 0, false
		}
		if !c.StatusIsActive("thoma-a1") {
			return 0, false
		}
		return float64(c.a1Stack) * 0.05, true
	})

	c.Core.Events.Subscribe(event.OnCharacterSwap, func(_ ...interface{}) bool {
		c.a1Stack = 0
		return false
	}, "thoma-a1-swap")
}
