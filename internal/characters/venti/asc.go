package venti

import "github.com/genshinsim/gcsim/pkg/core/attributes"

func (c *char) a4() {
	c.AddEnergy("venti-a4", 15)
	if c.qInfuse == attributes.NoElement {
		return
	}
	for _, char := range c.Core.Player.Chars() {
		if char.Base.Element == c.qInfuse {
			char.AddEnergy("venti-a4", 15)
		}
	}
}
