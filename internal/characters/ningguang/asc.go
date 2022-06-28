package ningguang

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/construct"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// activate a4 if screen is on-field and character uses dash
func (c *char) a4() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.GeoP] = 0.12
	//TODO: this used to be on PostDash; need to check if working correctly still
	c.Core.Events.Subscribe(event.OnDash, func(args ...interface{}) bool {
		// check for jade screen
		if c.Core.Constructs.CountByType(construct.GeoConstructNingSkill) <= 0 {
			return false
		}
		active := c.Core.Player.ActiveChar()
		active.AddStatMod(character.StatMod{Base: modifier.NewBase("ning-screen", 600), AffectedStat: attributes.GeoP, Amount: func() ([]float64, bool) {
			return m, true
		}})
		return false
	}, "ningguang-a4")
}
