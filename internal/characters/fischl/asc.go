package fischl

import "github.com/genshinsim/gsim/pkg/def"

func (c *char) a4() {
	last := 0
	c.Sim.AddOnReaction(func(t def.Target, ds *def.Snapshot) {
		if ds.ActorIndex != c.Sim.ActiveCharIndex() {
			return
		}
		//check reaction type, only care for overload, electro charge, superconduct
		switch ds.ReactionType {
		case def.Overload:
		case def.ElectroCharged:
		case def.Superconduct:
		case def.SwirlElectro:
		default:
			return
		}
		//do nothing if oz not on field
		if c.ozActiveUntil < c.Sim.Frame() {
			return
		}
		if c.Sim.Frame()-30 < last && last != 0 {
			return
		}
		last = c.Sim.Frame()

		d := c.Snapshot(
			"Fischl A4",
			def.AttackTagElementalArt,
			def.ICDTagNone,
			def.ICDGroupFischl,
			def.StrikeTypePierce,
			def.Electro,
			25,
			0.8,
		)
		d.Targets = t.Index()
		c.QueueDmg(&d, 1)

	}, "fischl a4")
}
