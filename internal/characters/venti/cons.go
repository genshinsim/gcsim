package venti

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

// TODO: the airborne part isn't implemented
// Skyward Sonnet decreases opponents' Anemo RES and Physical RES by 12% for 10s.
// Opponents launched by Skyward Sonnet suffer an additional 12% Anemo RES and Physical RES decrease while airborne.
func (c *char) c2(a combat.AttackCB) {
	if c.Base.Cons < 2 {
		return
	}
	e, ok := a.Target.(core.Enemy)
	if !ok {
		return
	}

	e.AddResistMod("venti-c2-anemo", 600, attributes.Anemo, -0.12)
	e.AddResistMod("venti-c2-phys", 600, attributes.Physical, -0.12)
}

// Targets who take DMG from Wind's Grand Ode have their Anemo RES decreased by 20%.
// If an Elemental Absorption occurred, then their RES towards the corresponding Element is also decreased by 20%.
func (c *char) c6(ele attributes.Element) func(a combat.AttackCB) {
	return func(a combat.AttackCB) {
		e, ok := a.Target.(core.Enemy)
		if !ok {
			return
		}
		e.AddResistMod("venti-c6-"+ele.String(), 600, ele, -0.20)
	}
}
