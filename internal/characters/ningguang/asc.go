package ningguang

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/construct"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// A1 is implemented in ningguang.go:
// When Ningguang is in possession of Star Jades, her Charged Attack does not consume Stamina.

// A character that passes through the Jade Screen will gain a 12% Geo DMG Bonus for 10s.
//
// - activate if screen is on-field and character uses dash
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.GeoP] = 0.12
	//TODO: this used to be on PostDash; need to check if working correctly still
	c.Core.Events.Subscribe(event.OnDash, func(_ ...interface{}) bool {
		// check for jade screen
		if c.Core.Constructs.CountByType(construct.GeoConstructNingSkill) <= 0 {
			return false
		}
		active := c.Core.Player.ActiveChar()
		active.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("ning-screen", 600),
			AffectedStat: attributes.GeoP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
		return false
	}, "ningguang-a4")
}
