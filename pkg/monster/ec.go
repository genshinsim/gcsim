package monster

import (
	"log"

	"github.com/genshinsim/gsim/pkg/def"
)

type AuraEC struct {
	*Element
	//EC will always carry two items, electro and hydro
	//any element applied will react with these 2 in order
	electro  *AuraElectro
	hydro    *AuraHydro
	t        *Target
	source   int
	snapshot def.Snapshot
}

func (a *AuraEC) AuraContains(ele ...def.EleType) bool {
	for _, v := range ele {
		if v == def.Electro || v == def.Hydro {
			return true
		}
	}
	return false
}

func newEC(e *AuraElectro, h *AuraHydro, t *Target, ds *def.Snapshot, f int) Aura {
	ec := AuraEC{}
	ec.Element = &Element{}
	ec.T = def.EC
	ec.electro = e
	ec.hydro = h
	ec.source = f
	ec.snapshot = ds.Clone()
	ec.t = t
	//on add, trigger tick immediately
	t.queueReaction(ds, def.ElectroCharged, 0, 1) //residual durability doesn't matter for EC
	//add handler to wane on dmg
	t.AddOnAttackLandedHook(
		func(ds *def.Snapshot) {
			if ds.AttackTag != def.AttackTagECDamage {
				return
			}
			//check if ec still active
			if t.aura.Type() != def.EC {
				return
			}
			v, ok := t.aura.(*AuraEC)
			if !ok {
				log.Panic("unexpected aura not type EC")
			}
			//wane in 0.1 seconds
			t.addTask(func(t *Target) {
				v.wane()
			}, 6)
		},
		"ec",
	)
	//add self repeating ticks
	t.addTask(ec.nextTick(f), 60)

	return &ec
}

func (a *AuraEC) nextTick(src int) func(t *Target) {
	return func(t *Target) {
		if a.source != src {
			//source changed, do nothing
			return
		}
		//ec SHOULD be active still, since if not we would have
		//called cleanup and set source to -1
		if a.t.aura != nil && a.t.aura.Type() != def.EC {
			return //redundant check anyways
		}
		//so ec is active, which means both aura must still have value > 0; so we can do dmg
		ds := a.snapshot.Clone()
		t.queueReaction(&ds, def.ElectroCharged, 0, 1)
		//queue up next tick
		a.t.addTask(a.nextTick(src), 60)
	}
}

func (a *AuraEC) cleanup(e, h bool) {
	if e || h {
		//both electro and hydro gone
		if e && h {
			a.t.aura = nil //nothing left
		} else if e { //just electro gone
			a.t.aura = a.hydro
		} else { //just hydro gone
			a.t.aura = a.electro
		}
		//cleanup
		a.t.RemoveOnAttackLandedHook("ec")
		a.source = -1
	}
}

func (a *AuraEC) wane() {
	//reduce both electro and hydro by 0.4
	//if either < 0 then we need to clean up ec
	a.electro.CurrentDurability -= 10
	a.hydro.CurrentDurability -= 10
	e := a.electro.CurrentDurability <= 0
	h := a.hydro.CurrentDurability <= 0
	a.cleanup(e, h)
}

func (a *AuraEC) Tick() bool {
	e := a.electro.Tick()
	h := a.hydro.Tick()

	a.cleanup(e, h)

	//otherwise all is well
	return false
}

func (a *AuraEC) React(ds *def.Snapshot, t *Target) (Aura, bool) {
	if ds.Durability == 0 {
		return a, false
	}
	switch ds.Element {
	case def.Anemo:
		//swirl electro first
		ds.ReactionType = def.SwirlElectro
		t.queueReaction(ds, def.SwirlElectro, a.electro.CurrentDurability, 1)
		a.electro.Reduce(ds, 0.5)

		if ds.Durability > 0 {
			ds.ReactionType = def.SwirlHydro
			//queue swirl dmg
			t.queueReaction(ds, def.SwirlHydro, a.hydro.CurrentDurability, 1)
			//reduce hydro by 0.5 of anemo
			a.hydro.Reduce(ds, 0.5)
		}
	case def.Geo:
		//for now assuming only crystallize electro
		ds.ReactionType = def.CrystallizeElectro
		shd := NewCrystallizeShield(def.Electro, t.sim.Frame(), ds.CharLvl, ds.Stats[def.EM], t.sim.Frame()+900)
		t.sim.AddShield(shd)
		//reduce by .05
		a.Reduce(ds, 0.5)
	case def.Pyro:
		//overload then vaporize
		ds.ReactionType = def.Overload
		t.queueReaction(ds, def.Overload, 0, 1)
		a.electro.Reduce(ds, 1)

		if ds.Durability > 0 {
			ds.ReactionType = def.Vaporize
			ds.ReactMult = 1.5
			a.hydro.Reduce(ds, 0.5)
		}
	case def.Hydro:
		//refresh hydro, update snapshot, and trigger 1 tick
		a.hydro.Refresh(ds.Durability)
		a.snapshot = ds.Clone()
		a.source = t.sim.Frame()
		//trigger tick and update tick timer
		t.queueReaction(ds, def.ElectroCharged, 0, 1)
		t.addTask(a.nextTick(t.sim.Frame()), 60)
		ds.ReactionType = def.ElectroCharged
	case def.Cryo:
		//superconduct and if any left trigger freeze
		ds.ReactionType = def.Superconduct
		t.queueReaction(ds, def.Superconduct, 0, 1)
		a.electro.Reduce(ds, 1)

		if ds.Durability > 0 {
			log.Println("FREEZE ON EC NOT IMPLEMENTED")
		}
	case def.Electro:
		//refresh electro, update snapshot, and trigger 1 tick
		a.electro.Refresh(ds.Durability)
		a.snapshot = ds.Clone()
		a.source = t.sim.Frame()
		//trigger tick and update tick timer
		t.queueReaction(ds, def.ElectroCharged, 0, 1)
		t.addTask(a.nextTick(t.sim.Frame()), 60)
		ds.ReactionType = def.ElectroCharged
	default:
		return a, false
	}
	e := a.electro.CurrentDurability <= 0
	h := a.hydro.CurrentDurability <= 0

	if e || h {
		a.cleanup(e, h)
		return t.aura, true
	}
	return a, true
}
