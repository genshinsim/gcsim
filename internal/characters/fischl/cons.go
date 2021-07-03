package fischl

import "github.com/genshinsim/gsim/pkg/def"

func (c *char) c6() {
	c.Sim.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		//do nothing if oz not on field
		if c.ozActiveUntil < c.Sim.Frame() {
			return
		}
		switch ds.AttackTag {
		case def.AttackTagNormal:
		case def.AttackTagTartagliaAttack:
		case def.AttackTagGandalfrAttack:
		default:
			return
		}

		d := c.Snapshot(
			"Fischl C6",
			def.AttackTagElementalArt,
			def.ICDTagElementalArt,
			def.ICDGroupFischl,
			def.StrikeTypePierce,
			def.Electro,
			25,
			0.3,
		)
		d.Targets = t.Index()
		c.QueueDmg(&d, 1)
	}, "fischl c6")
}
