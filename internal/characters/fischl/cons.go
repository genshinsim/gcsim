package fischl

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) c6() {
	//this is on attack animation state, not attack landed
	c.Core.Events.Subscribe(core.PostAttack, func(args ...interface{}) bool {
		//do nothing if oz not on field
		if c.ozActiveUntil < c.Core.F {
			return false
		}

		d := c.Snapshot(
			"Fischl C6",
			core.AttackTagElementalArt,
			core.ICDTagElementalArt,
			core.ICDGroupFischl,
			core.StrikeTypePierce,
			core.Electro,
			25,
			0.3,
		)
		d.Targets = 0
		c.QueueDmg(&d, 1)
		return false
	}, "fischl-c6")
}
