package raiden

import (
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

// When nearby party members gain Elemental Orbs or Particles, Chakra Desiderata gains 2 Resolve stacks.
// This effect can occur once every 3s.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	particleICD := 0
	c.Core.Events.Subscribe(event.OnParticleReceived, func(_ ...interface{}) bool {
		if particleICD > c.Core.F {
			return false
		}
		particleICD = c.Core.F + 180 // once every 3 seconds
		c.stacks += 2
		if c.stacks > 60 {
			c.stacks = 60
		}
		return false
	}, "raiden-particle-stacks")
}

// Each 1% above 100% Energy Recharge that the Raiden Shogun possesses grants her:
//
// - 0.6% greater Energy restoration from Musou Isshin
func (c *char) a4Energy(er float64) float64 {
	if c.Base.Ascension < 4 {
		return 0
	}
	excess := int(er / 0.01)
	increase := float64(excess) * 0.006
	c.Core.Log.NewEvent("a4 energy restore stacks", glog.LogCharacterEvent, c.Index).
		Write("stacks", excess).
		Write("increase", increase)
	return increase
}

// This is implemented in raiden.go:
// Each 1% above 100% Energy Recharge that the Raiden Shogun possesses grants her:
//
// - 0.4% Electro DMG Bonus.
