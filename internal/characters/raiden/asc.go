package raiden

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
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
	excess := int(er * 100)
	increase := float64(excess) * 0.006
	c.Core.Log.NewEvent("a4 energy restore stacks", glog.LogCharacterEvent, c.Index).
		Write("stacks", excess).
		Write("increase", increase)
	return increase
}

// Each 1% above 100% Energy Recharge that the Raiden Shogun possesses grants her:
//
// - 0.4% Electro DMG Bonus.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	if c.a4Stats == nil {
		c.a4Stats = make([]float64, attributes.EndStatType)
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("raiden-a4", -1),
			AffectedStat: attributes.ElectroP,
			Extra:        true,
			Amount: func() ([]float64, bool) {
				return c.a4Stats, true
			},
		})
	}
	c.a4Stats[attributes.ElectroP] = c.NonExtraStat(attributes.ER) * 0.4 // 100 * 0.004
	c.QueueCharTask(c.a4, 30)
}
