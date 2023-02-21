package venti

import "github.com/genshinsim/gcsim/pkg/core/attributes"

// A1 is not implemented and will likely never be implemented:
// Holding Skyward Sonnet creates an upcurrent that lasts for 20s.

// Regenerates 15 Energy for Venti after the effects of Wind's Grand Ode end.
// If an Elemental Absorption occurred, this also restores 15 Energy to all characters of that corresponding element in the party.
//
// - checks for ascension level in burst.go to avoid queuing this up only to fail the ascension level check
func (c *char) a4() {
	c.AddEnergy("venti-a4", 15)
	if c.qAbsorb == attributes.NoElement {
		return
	}
	for _, char := range c.Core.Player.Chars() {
		if char.Base.Element == c.qAbsorb {
			char.AddEnergy("venti-a4", 15)
		}
	}
}
