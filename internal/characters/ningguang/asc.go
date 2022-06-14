package ningguang

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/construct"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

// activate a4 if screen is on-field and character uses dash
func (c *char) a4() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.GeoP] = 0.12
	c.Core.Events.Subscribe(event.PostDash, func(args ...interface{}) bool {
		// check for jade screen
		if c.Core.Constructs.CountByType(construct.GeoConstructNingSkill) <= 0 {
			return false
		}
		active := c.Core.Player.ActiveChar()
		active.AddStatMod("ning-screen", 600, attributes.GeoP, func() ([]float64, bool) {
			return m, true
		})
		return false
	}, "ningguang-a4")
}
