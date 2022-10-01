package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

func calcSwirlAtkDurability(consumed, src combat.Durability) combat.Durability {
	if consumed < src {
		return 1.25*(0.5*consumed-1) + 25
	}
	return 1.25*(src-1) + 25
}

func (r *Reactable) queueSwirl(rt combat.ReactionType, ele attributes.Element, tag combat.AttackTag, icd combat.ICDTag, dur combat.Durability, charIndex int) {
	//swirl triggers two attacks; one self with no gauge
	//and one aoe with gauge
	ai := combat.AttackInfo{
		ActorIndex:       charIndex,
		DamageSrc:        r.self.Key(),
		Abil:             string(rt),
		AttackTag:        tag,
		ICDTag:           icd,
		ICDGroup:         combat.ICDGroupReactionA,
		Element:          ele,
		IgnoreDefPercent: 1,
	}
	char := r.core.Player.ByIndex(charIndex)
	em := char.Stat(attributes.EM)
	ai.FlatDmg = 0.6 * calcReactionDmg(char, ai, em)
	snap := combat.Snapshot{
		CharLvl:  char.Base.Level,
		ActorEle: char.Base.Element,
	}
	snap.Stats[attributes.EM] = em
	//first attack is self no hitbox
	r.core.QueueAttackWithSnap(
		ai,
		snap,
		combat.NewDefSingleTarget(r.self.Index(), r.self.Type()),
		1,
	)
	//next is aoe - hydro swirls never do AoE damage, as they only spread the element
	if ele == attributes.Hydro {
		ai.FlatDmg = 0
	}
	ai.Durability = dur
	ai.Abil = string(rt) + " (aoe)"
	r.core.QueueAttackWithSnap(
		ai,
		snap,
		combat.NewCircleHit(r.self, 5, false, combat.TargettableEnemy, combat.TargettableGadget),
		1,
	)
}

func (r *Reactable) trySwirlElectro(a *combat.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}
	if r.Durability[ModifierElectro] < ZeroDur {
		return
	}
	rd := r.reduce(attributes.Electro, a.Info.Durability, 0.5)
	atkDur := calcSwirlAtkDurability(rd, a.Info.Durability)
	a.Info.Durability -= rd
	a.Reacted = true
	//queue an attack first
	r.core.Events.Emit(event.OnSwirlElectro, r.self, a)
	r.queueSwirl(
		combat.SwirlElectro,
		attributes.Electro,
		combat.AttackTagSwirlElectro,
		combat.ICDTagSwirlElectro,
		atkDur,
		a.Info.ActorIndex,
	)
	//at this point if any durability left, we need to check for prescence of
	//hydro in case of EC
	if a.Info.Durability > ZeroDur && r.Durability[ModifierHydro] > ZeroDur {
		//trigger swirl hydro
		r.trySwirlHydro(a)
		//check EC clean up
		r.checkEC()
	}
}

func (r *Reactable) trySwirlHydro(a *combat.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}
	if r.Durability[ModifierHydro] < ZeroDur {
		return
	}
	rd := r.reduce(attributes.Hydro, a.Info.Durability, 0.5)
	atkDur := calcSwirlAtkDurability(rd, a.Info.Durability)
	a.Info.Durability -= rd
	a.Reacted = true
	//queue an attack first
	r.core.Events.Emit(event.OnSwirlHydro, r.self, a)
	r.queueSwirl(
		combat.SwirlHydro,
		attributes.Hydro,
		combat.AttackTagSwirlHydro,
		combat.ICDTagSwirlHydro,
		atkDur,
		a.Info.ActorIndex,
	)
}

func (r *Reactable) trySwirlCryo(a *combat.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}
	if r.Durability[ModifierCryo] < ZeroDur {
		return
	}
	rd := r.reduce(attributes.Cryo, a.Info.Durability, 0.5)
	atkDur := calcSwirlAtkDurability(rd, a.Info.Durability)
	a.Info.Durability -= rd
	a.Reacted = true
	//queue an attack first
	r.core.Events.Emit(event.OnSwirlCryo, r.self, a)
	r.queueSwirl(
		combat.SwirlCryo,
		attributes.Cryo,
		combat.AttackTagSwirlCryo,
		combat.ICDTagSwirlCryo,
		atkDur,
		a.Info.ActorIndex,
	)
}

func (r *Reactable) trySwirlPyro(a *combat.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}
	if r.Durability[ModifierPyro] < ZeroDur {
		return
	}
	rd := r.reduce(attributes.Pyro, a.Info.Durability, 0.5)
	atkDur := calcSwirlAtkDurability(rd, a.Info.Durability)
	a.Info.Durability -= rd
	a.Reacted = true
	r.burningCheck()
	//queue an attack first
	r.core.Events.Emit(event.OnSwirlPyro, r.self, a)
	r.queueSwirl(
		combat.SwirlPyro,
		attributes.Pyro,
		combat.AttackTagSwirlPyro,
		combat.ICDTagSwirlPyro,
		atkDur,
		a.Info.ActorIndex,
	)
}

func (r *Reactable) trySwirlFrozen(a *combat.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}
	if r.Durability[ModifierFrozen] < ZeroDur {
		return
	}
	rd := r.reduce(attributes.Frozen, a.Info.Durability, 0.5)
	atkDur := calcSwirlAtkDurability(rd, a.Info.Durability)
	a.Info.Durability -= rd
	a.Reacted = true
	//queue an attack first
	r.core.Events.Emit(event.OnSwirlCryo, r.self, a)
	r.queueSwirl(
		combat.SwirlCryo,
		attributes.Cryo,
		combat.AttackTagSwirlCryo,
		combat.ICDTagSwirlCryo,
		atkDur,
		a.Info.ActorIndex,
	)
	r.checkFreeze()
}
