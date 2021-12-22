package reactable

import "github.com/genshinsim/gcsim/pkg/core"

func (r *Reactable) tryMelt(a *core.AttackEvent) {
	if a.Info.Durability < zeroDur {
		return
	}
	if r.Durability[core.Frozen] > zeroDur {
		return
	}
	switch a.Info.Element {
	case core.Pyro:
		if r.Durability[core.Cryo] < zeroDur {
			return
		}
		r.reduce(core.Cryo, a.Info.Durability, 2)
		a.Info.AmpMult = 2.0
	case core.Cryo:
		if r.Durability[core.Pyro] < zeroDur {
			return
		}
		r.reduce(core.Pyro, a.Info.Durability, 0.5)
		a.Info.AmpMult = 1.5
	default:
		//should be here
		return
	}
	//there shouldn't be anything else to react with if not frozen
	a.Info.Durability = 0
	a.Info.Amped = true
	a.Info.AmpType = core.Melt
	r.core.Events.Emit(core.OnMelt, r.self, a)
}

func (r *Reactable) tryMeltFrozen(a *core.AttackEvent) {
	if a.Info.Durability < zeroDur {
		return
	}
	if r.Durability[core.Frozen] < zeroDur {
		return
	}
	switch a.Info.Element {
	case core.Pyro:
		//TODO: the assumption here is we first reduce cryo, and if there's any
		//src durability left, we reduce frozen. note that it's still only one
		//melt reaction
		a.Info.Durability -= r.reduce(core.Cryo, a.Info.Durability, 2)
		r.reduce(core.Frozen, a.Info.Durability, 2)
		a.Info.AmpMult = 2.0
	default:
		//should be here
		return
	}
	//durability not wiped out because we can potentially react with hydro still
	a.Info.Amped = true
	a.Info.AmpType = core.Melt
	r.core.Events.Emit(core.OnMelt, r.self, a)
}
