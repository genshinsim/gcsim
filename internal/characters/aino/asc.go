package aino

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

func (c *char) a1BurstEnhance() (int, float64, attacks.ICDGroup, attacks.ICDTag) {
	if c.Base.Ascension < 1 {
		return 90, 1, attacks.ICDGroupDefault, attacks.ICDTagElementalBurst
	}

	if c.Core.Player.GetMoonsignLevel() < 2 {
		return 90, 1, attacks.ICDGroupDefault, attacks.ICDTagElementalBurst
	}
	return 42, 4, attacks.ICDGroupAinoBurstMoonHit, attacks.ICDTagAinoBurstMoonHit
}

func (c *char) a4Dmg() float64 {
	if c.Base.Ascension < 4 {
		return 0
	}
	return c.Stat(attributes.EM) * 0.5
}
