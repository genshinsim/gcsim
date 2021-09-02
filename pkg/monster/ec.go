package monster

import (
	"fmt"
	"log"

	"github.com/genshinsim/gsim/pkg/core"
)

type AuraEC struct {
	*Element
	//EC will always carry two items, electro and hydro
	//any element applied will react with these 2 in order
	electro  *AuraElectro
	hydro    *AuraHydro
	t        *Target
	source   int
	snapshot core.Snapshot
}

func (a *AuraEC) AuraContains(ele ...core.EleType) bool {
	for _, v := range ele {
		if v == core.Electro || v == core.Hydro {
			return true
		}
	}
	return false
}

func newEC(e *AuraElectro, h *AuraHydro, t *Target, ds *core.Snapshot, f int) Aura {
	ec := AuraEC{}
	ec.Element = &Element{}
	ec.T = core.EC
	ec.electro = e
	ec.hydro = h
	ec.source = f
	ec.snapshot = ds.Clone()
	ec.t = t
	//on add, trigger tick in 0.1s
	t.core.Tasks.Add(func() {
		t.queueReaction(ds, core.ElectroCharged, 0, 1) //residual durability doesn't matter for EC
	}, 10)
	//add handler to wane on dmg; should ignore if dmg = 0
	t.core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		//target should be first, then snapshot
		n := args[0].(core.Target)
		snap := args[1].(*core.Snapshot)
		dmg := args[2].(float64)
		if n.Index() != t.index {
			return false
		}
		if snap.AttackTag != core.AttackTagECDamage {
			return false
		}
		if t.aura == nil {
			return true
		}
		//check if ec still active
		if t.aura.Type() != core.EC {
			return true
		}
		//ignore if this dmg instance has been wiped out due to icd
		if dmg == 0 {
			return false
		}
		v, ok := t.aura.(*AuraEC)
		if !ok {
			log.Panic("unexpected aura not type EC")
		}
		//wane in 0.1 seconds
		t.core.Tasks.Add(func() {
			v.wane()
		}, 6)
		return false
	}, fmt.Sprintf("ec-%v", t.index))
	t.core.Tasks.Add(ec.nextTick(f), 60+10)
	//add self repeating ticks

	return &ec
}

func (a *AuraEC) nextTick(src int) func() {
	return func() {
		if a.source != src {
			//source changed, do nothing
			return
		}
		//ec SHOULD be active still, since if not we would have
		//called cleanup and set source to -1
		if a.t.aura != nil && a.t.aura.Type() != core.EC {
			return //redundant check anyways
		}
		//so ec is active, which means both aura must still have value > 0; so we can do dmg
		ds := a.snapshot.Clone()
		a.t.queueReaction(&ds, core.ElectroCharged, 0, 1)
		//queue up next tick
		a.t.core.Tasks.Add(a.nextTick(src), 60)
	}
}

func (a *AuraEC) cleanup(e, h bool) {
	if e || h {
		//both electro and hydro gone
		if e && h {
			a.t.core.Log.Debugw("ec expired",
				"frame", a.t.core.F,
				"event", core.LogElementEvent,
				"aura", "ec",
				"target", a.t.index,
			)
			a.t.aura = nil //nothing left

		} else if e { //just electro gone
			a.t.core.Log.Debugw("electro (in ec) expired",
				"frame", a.t.core.F,
				"event", core.LogElementEvent,
				"aura", "ec",
				"target", a.t.index,
			)
			a.t.aura = a.hydro
		} else { //just hydro gone
			a.t.core.Log.Debugw("hydro (in ec) expired",
				"frame", a.t.core.F,
				"event", core.LogElementEvent,
				"aura", "ec",
				"target", a.t.index,
			)
			a.t.aura = a.electro
		}
		//cleanup
		a.t.core.Events.Unsubscribe(core.OnDamage, fmt.Sprintf("ec-%v", a.t.index))
		a.source = -1
	}
}

func (a *AuraEC) wane() {
	//reduce both electro and hydro by 0.4
	//if either < 0 then we need to clean up ec
	a.electro.CurrentDurability -= 10
	a.hydro.CurrentDurability -= 10
	a.t.core.Log.Debugw("ec wane",
		"frame", a.t.core.F,
		"event", core.LogElementEvent,
		"aura", "ec",
		"electro", a.electro.CurrentDurability,
		"hydro", a.hydro.CurrentDurability,
		"target", a.t.index,
	)

	e := a.electro.CurrentDurability <= 0
	h := a.hydro.CurrentDurability <= 0
	a.cleanup(e, h)
}

func (a *AuraEC) Tick() bool {
	e := a.electro.Tick()
	h := a.hydro.Tick()

	//TODO: CHECK FOR JUMPING

	a.cleanup(e, h)

	//otherwise all is well
	return false
}

func (a *AuraEC) React(ds *core.Snapshot, t *Target) (Aura, bool) {
	if ds.Durability == 0 {
		return a, false
	}
	switch ds.Element {
	case core.Anemo:
		//swirl electro first
		ds.ReactionType = core.SwirlElectro
		t.queueReaction(ds, core.SwirlElectro, a.electro.CurrentDurability, 1)
		a.electro.Reduce(ds, 0.5)

		if ds.Durability > 0 {
			ds.ReactionType = core.SwirlHydro
			//queue swirl dmg
			t.queueReaction(ds, core.SwirlHydro, a.hydro.CurrentDurability, 1)
			//reduce hydro by 0.5 of anemo
			a.hydro.Reduce(ds, 0.5)
		}
	case core.Geo:
		//for now assuming only crystallize electro
		ds.ReactionType = core.CrystallizeElectro
		shd := NewCrystallizeShield(core.Electro, t.core.F, ds.CharLvl, ds.Stats[core.EM], t.core.F+900)
		t.core.Shields.Add(shd)
		//reduce by .05
		a.Reduce(ds, 0.5)
	case core.Pyro:
		//overload then vaporize
		ds.ReactionType = core.Overload
		t.queueReaction(ds, core.Overload, 0, 1)
		a.electro.Reduce(ds, 1)

		if ds.Durability > 0 {
			ds.ReactionType = core.Vaporize
			ds.ReactMult = 1.5
			ds.IsMeltVape = true
			a.hydro.Reduce(ds, 0.5)
		}
	case core.Hydro:
		//refresh hydro, update snapshot, and trigger 1 tick
		a.hydro.Refresh(ds.Durability)
		a.snapshot = ds.Clone()
		// a.source = t.core.F
		//trigger tick and update tick timer
		// t.queueReaction(ds, core.ElectroCharged, 0, 1)
		// t.core.Tasks.Add(a.nextTick(t.core.F), 60)
		ds.ReactionType = core.ElectroCharged
	case core.Cryo:
		//superconduct and if any left trigger freeze
		ds.ReactionType = core.Superconduct
		t.queueReaction(ds, core.Superconduct, 0, 1)
		a.electro.Reduce(ds, 1)

		if ds.Durability > 0 {
			log.Println("FREEZE ON EC NOT IMPLEMENTED")
		}
	case core.Electro:
		//refresh electro, update snapshot, and trigger 1 tick
		a.electro.Refresh(ds.Durability)
		a.snapshot = ds.Clone()
		// a.source = t.core.F
		//trigger tick and update tick timer
		// t.queueReaction(ds, core.ElectroCharged, 0, 1)
		// t.core.Tasks.Add(a.nextTick(t.core.F), 60)
		ds.ReactionType = core.ElectroCharged
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
