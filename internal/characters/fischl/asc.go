package fischl

import "github.com/genshinsim/gsim/pkg/core"

func (c *char) a4() {
	last := 0
	c.Core.Events.Subscribe(core.OnTransReaction, func(args ...interface{}) bool {
		t := args[0].(core.Target)
		ds := args[1].(*core.Snapshot)
		if ds.ActorIndex != c.Core.ActiveChar {
			return false
		}
		//check reaction type, only care for overload, electro charge, superconduct
		switch ds.ReactionType {
		case core.Overload:
		case core.ElectroCharged:
		case core.Superconduct:
		case core.SwirlElectro:
		default:
			return false
		}
		//do nothing if oz not on field
		if c.ozActiveUntil < c.Core.F {
			return false
		}
		if c.Core.F-30 < last && last != 0 {
			return false
		}
		last = c.Core.F

		d := c.Snapshot(
			"Fischl A4",
			core.AttackTagElementalArt,
			core.ICDTagNone,
			core.ICDGroupFischl,
			core.StrikeTypePierce,
			core.Electro,
			25,
			0.8,
		)
		d.Targets = t.Index()
		c.QueueDmg(&d, 1)

		return false
	}, "fischl-a4")

}
