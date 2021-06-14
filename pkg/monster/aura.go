package monster

import "github.com/genshinsim/gsim/pkg/def"

type Aura interface {
	React(ds *def.Snapshot, s def.Sim, auras []Aura) bool //react and modify the slice of auras, return true if reacted
	Tick()                                                //tick down the aura
	Durability() def.Durability                           //return the current durability
	Attach(ele def.EleType, dur def.Durability, f int)    //attach aura, refresh if existing
	Source() int                                          //return the origination source of this aura, used to distinguish if diff
}

func AttachAura(ds *def.Snapshot) Aura {
	var a Aura
	switch ds.Element {
	case def.Pyro:
		r := &PyroAura{}
		r.Element = &Element{}
		a = r
	case def.Hydro:
		r := &HydroAura{}
		r.Element = &Element{}
		a = r
	case def.Cryo:
		r := &CryoAura{}
		r.Element = &Element{}
		a = r
	case def.Electro:
		r := &ElectroAura{}
		r.Element = &Element{}
		a = r
	default:
		return nil

	}
	return a
}

type Element struct {
	Type              def.EleType
	MaxDurability     def.Durability
	CurrentDurability def.Durability
	DecayRate         def.Durability //amount of durability decay per tick
	Start             int
}

func (e *Element) React(ds *def.Snapshot, s def.Sim, auras []Aura) bool {
	return false
}

func (e *Element) Tick() {
	e.CurrentDurability -= e.DecayRate
}

func (e *Element) Durability() def.Durability {
	return e.CurrentDurability
}

func (e *Element) Attach(ele def.EleType, dur def.Durability, f int) {
	e.Start = f
	e.Type = ele
	e.MaxDurability = 0.8 * dur
	e.CurrentDurability = 0.8 * dur

	//duration = 0.1 * dur + 7
	//in frames that's 6 * dur + 420
	//rate is therefore dur / (6 * dur + 420)
	e.DecayRate = dur / (6*dur + 420)
}

func (e *Element) Refresh(dur def.Durability) {
	e.CurrentDurability += dur
}

func (e *Element) Source() int {
	return e.Start
}
