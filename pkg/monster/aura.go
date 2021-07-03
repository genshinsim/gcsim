package monster

import (
	"github.com/genshinsim/gsim/pkg/def"
)

type Aura interface {
	React(ds *def.Snapshot, t *Target) (Aura, bool) //react and modify the slice of auras, return true if reacted
	Tick() bool                                     //tick down the aura
	Durability() def.Durability                     //return the current durability
	Attach(dur def.Durability, f int)               //attach aura, refresh if existing
	Source() int                                    //return the origination source of this aura, used to distinguish if diff
	Type() def.EleType
	AuraContains(e ...def.EleType) bool
}

func NewAura(ds *def.Snapshot, f int) Aura {
	var a Aura
	switch ds.Element {
	case def.Pyro:
		r := &AuraPyro{}
		r.Element = &Element{}
		r.T = def.Pyro
		a = r
	case def.Hydro:
		r := &AuraHydro{}
		r.Element = &Element{}
		r.T = def.Hydro
		a = r
	case def.Cryo:
		r := &AuraCyro{}
		r.Element = &Element{}
		r.T = def.Cryo
		a = r
	case def.Electro:
		r := &AuraElectro{}
		r.Element = &Element{}
		r.T = def.Electro
		a = r
	default:
		return nil

	}
	a.Attach(ds.Durability, f)
	return a
}

type Element struct {
	T                 def.EleType
	MaxDurability     def.Durability
	CurrentDurability def.Durability
	DecayRate         def.Durability //amount of durability decay per tick
	Start             int
}

func (e *Element) AuraContains(ele ...def.EleType) bool {
	for _, v := range ele {
		if v == e.T {
			return true
		}
	}
	return false
}

func (e *Element) Type() def.EleType {
	return e.T
}

func (e *Element) Tick() bool {
	e.CurrentDurability -= e.DecayRate
	return e.CurrentDurability < 0
}

func (e *Element) Durability() def.Durability {
	return e.CurrentDurability
}

func (e *Element) Attach(dur def.Durability, f int) {
	e.Start = f
	e.MaxDurability = 0.8 * dur
	e.CurrentDurability = 0.8 * dur

	//duration = 0.1 * dur + 7
	//in frames that's 6 * dur + 420
	//rate is therefore dur / (6 * dur + 420)
	e.DecayRate = 0.8 * dur / (6*dur + 420)
}

func (e *Element) Refresh(dur def.Durability) {
	e.CurrentDurability += dur
}

func (e *Element) Source() int {
	return e.Start
}

func (e *Element) Reduce(ds *def.Snapshot, factor def.Durability) def.Durability {
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
