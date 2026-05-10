package aino

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

func (c *char) a1BurstEnhance() (int, float64, attacks.ICDGroup, attacks.ICDTag) {
	if c.Base.Ascension < 1 || c.Core.Player.GetMoonsignLevel() < 2 {
		return 1.5 * 60, 1, attacks.ICDGroupDefault, attacks.ICDTagElementalBurst
	}
	return 0.7 * 60, 4, attacks.ICDGroupAinoBurstMoonHit, attacks.ICDTagAinoBurstMoonHit
}

func (c *char) a4Dmg() float64 {
	if c.Base.Ascension < 4 {
		return 0
	}
	return c.Stat(attributes.EM) * 0.5
}
