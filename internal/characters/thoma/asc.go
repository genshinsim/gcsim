package thoma

import (
	"github.com/genshinsim/gcsim/pkg/core/event"
)

// When your current active character obtains or refreshes a Blazing Barrier,
// this character's Shield Strength will increase by 5% for 6s.
// This effect can be triggered once every 0.3s seconds. Max 5 stacks.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
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

// DMG dealt by Crimson Ooyoroi's Fiery Collapse is increased by 2.2% of Thoma's Max HP.
func (c *char) a4() float64 {
	if c.Base.Ascension < 1 {
		return 0
	}
	return 0.022 * c.MaxHP()
}
