package fischl

import "github.com/genshinsim/gsim/pkg/core"

func (c *char) c6() {
	//this is on attack animation state, not attack landed
	c.Sim.AddEventHook(func(s core.Sim) bool {
		//do nothing if oz not on field
		if c.ozActiveUntil < c.Sim.Frame() {
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
	}, "fischl c6", core.PostAttackHook)

	// c.Sim.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
	// 	//do nothing if oz not on field
	// 	if c.ozActiveUntil < c.Sim.Frame() {
	// 		return
	// 	}
	// 	switch ds.AttackTag {
	// 	case def.AttackTagNormal:
	// 	case def.AttackTagTartagliaAttack:
	// 	case def.AttackTagGandalfrAttack:
	// 	default:
	// 		return
	// 	}

	// 	d := c.Snapshot(
	// 		"Fischl C6",
	// 		def.AttackTagElementalArt,
	// 		def.ICDTagElementalArt,
	// 		def.ICDGroupFischl,
	// 		def.StrikeTypePierce,
	// 		def.Electro,
	// 		25,
	// 		0.3,
	// 	)
	// 	d.Targets = t.Index()
	// 	c.QueueDmg(&d, 1)
	// }, "fischl c6")
}
