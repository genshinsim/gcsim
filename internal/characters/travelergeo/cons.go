package travelergeo

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/construct"
)

//Party members within the radius of Wake of Earth have their CRIT Rate increased by 10%
//and have increased resistance against interruption.
func (c *char) c1() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = .1
	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod("geo-traveler-c1", -1, attributes.CR, func() ([]float64, bool) {
			if c.Core.Constructs.CountByType(construct.GeoConstructTravellerBurst) == 0 {
				return nil, false
			}
			return m, true
		})
	}
}
