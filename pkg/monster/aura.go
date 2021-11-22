package monster

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

type Aura interface {
	React(ds *core.Snapshot, t *Target) (Aura, bool) //react and modify the slice of auras, return true if reacted
	Tick() bool                                      //tick down the aura
	Durability() core.Durability                     //return the current durability
	Attach(dur core.Durability, f int)               //attach aura, refresh if existing
	Source() int                                     //return the origination source of this aura, used to distinguish if diff
	Type() core.EleType
	AuraContains(e ...core.EleType) bool
}

func NewAura(ds *core.Snapshot, f int) Aura {
	var a Aura
	switch ds.Element {
	case core.Pyro:
		r := &AuraPyro{}
		r.Element = &Element{}
		r.T = core.Pyro
		a = r
	case core.Hydro:
		r := &AuraHydro{}
		r.Element = &Element{}
		r.T = core.Hydro
		a = r
	case core.Cryo:
		r := &AuraCyro{}
		r.Element = &Element{}
		r.T = core.Cryo
		a = r
	case core.Electro:
		r := &AuraElectro{}
		r.Element = &Element{}
		r.T = core.Electro
		a = r
	default:
		return nil

	}
	a.Attach(ds.Durability, f)
	return a
}

type Element struct {
	T                 core.EleType
	MaxDurability     core.Durability
	CurrentDurability core.Durability
	DecayRate         core.Durability //amount of durability decay per tick
	Start             int
}

func (e *Element) AuraContains(ele ...core.EleType) bool {
	for _, v := range ele {
		if v == e.T {
			return true
		}
	}
	return false
}

func (e *Element) Type() core.EleType {
	return e.T
}

func (e *Element) Tick() bool {
	e.CurrentDurability -= e.DecayRate
	return e.CurrentDurability < 0
}

func (e *Element) Durability() core.Durability {
	return e.CurrentDurability
}

func (e *Element) Attach(dur core.Durability, f int) {
	e.Start = f
	e.MaxDurability = 0.8 * dur
	e.CurrentDurability = 0.8 * dur

	//duration = 0.1 * dur + 7
	//in frames that's 6 * dur + 420
	//rate is therefore dur / (6 * dur + 420)
	e.DecayRate = 0.8 * dur / (6*dur + 420)
}

func (e *Element) Refresh(dur core.Durability) {
	//refresh should only add 80%, subject to max of current
	e.CurrentDurability += 0.8 * dur
	if e.CurrentDurability > 0.8*dur {
		e.CurrentDurability = 0.8 * dur
	}
}

func (e *Element) Source() int {
	return e.Start
}

func (e *Element) Reduce(ds *core.Snapshot, factor core.Durability) core.Durability {
	//reduce current by the lower of current and factor * ds.Dur
	a := factor * ds.Durability
	b := ds.Durability
	if a > e.CurrentDurability {
		a = e.CurrentDurability
		b = a / factor
	}

	e.CurrentDurability -= a
	ds.Durability -= b

	return a
}
