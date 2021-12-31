package reactable

import "github.com/genshinsim/gcsim/pkg/core"

func (r *Reactable) tryVaporize(a *core.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}
	switch a.Info.Element {
	case core.Pyro:
		//make sure there's hydro
		if r.Durability[core.Hydro] < ZeroDur {
			return
		}
		//if there's still frozen left don't try to vape
		if r.Durability[core.Frozen] > ZeroDur {
			return
		}
		r.reduce(core.Hydro, a.Info.Durability, .5)
		a.Info.AmpMult = 1.5
	case core.Hydro:
		//make sure there's pyro to vape; no coexistance with pyro (yet)
		if r.Durability[core.Pyro] < ZeroDur {
			return
		}
		r.reduce(core.Pyro, a.Info.Durability, 2)
		a.Info.AmpMult = 2
	default:
		//should be here
		return
	}
	//there shouldn't be anything else to react with
	a.Info.Durability = 0
	a.Info.Amped = true
	a.Info.AmpType = core.Vaporize
	r.core.Events.Emit(core.OnVaporize, r.self, a)
}
