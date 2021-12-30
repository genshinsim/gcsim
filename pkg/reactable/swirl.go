package reactable

import "github.com/genshinsim/gcsim/pkg/core"

func calcSwirlAtkDurability(consumed, src core.Durability) core.Durability {
	if consumed < 0.5*src {
		return 1.25*(consumed-1) + 25
	}
	return 1.25*(src-1) + 25
}

func (r *Reactable) queueSwirl(rt core.ReactionType, ele core.EleType, tag core.AttackTag, icd core.ICDTag, dur core.Durability, charIndex int) {
	//swirl triggers two attacks; one self with no gauge
	//and one aoe with gauge
	ai := core.AttackInfo{
		ActorIndex:       charIndex,
		DamageSrc:        r.self.Index(),
		Abil:             string(rt),
		AttackTag:        tag,
		ICDTag:           icd,
		ICDGroup:         core.ICDGroupReactionA,
		Element:          ele,
		IgnoreDefPercent: 1,
	}
	ai.FlatDmg = 0.6 * r.calcReactionDmg(ai)
	//first attack is self no hitbox
	r.core.Combat.QueueAttack(
		ai,
		core.NewDefSingleTarget(r.self.Index(), r.self.Type()),
		-1,
		1,
	)
	//next is aoe
	ai.Durability = dur
	ai.Abil = string(rt) + " (aoe)"
	r.core.Combat.QueueAttack(
		ai,
		core.NewDefCircHit(5, false, core.TargettableEnemy),
		-1,
		1,
	)
}

func (r *Reactable) trySwirlElectro(a *core.AttackEvent) {
	if a.Info.Durability < zeroDur {
		return
	}
	if r.Durability[core.Electro] < zeroDur {
		return
	}
	rd := r.reduce(core.Electro, a.Info.Durability, 0.5)
	atkDur := calcSwirlAtkDurability(rd, a.Info.Durability)
	a.Info.Durability -= rd
	//queue an attack first
	r.core.Events.Emit(core.OnSwirlElectro, r.self, a)
	r.queueSwirl(
		core.SwirlElectro,
		core.Electro,
		core.AttackTagSwirlElectro,
		core.ICDTagSwirlElectro,
		atkDur,
		a.Info.ActorIndex,
	)
	//at this point if any durability left, we need to check for prescence of
	//hydro in case of EC
	if a.Info.Durability > zeroDur && r.Durability[core.Hydro] > zeroDur {
		//trigger swirl hydro
		r.trySwirlHydro(a)
		//check EC clean up
		r.checkEC()
	}
}

func (r *Reactable) trySwirlHydro(a *core.AttackEvent) {
	if a.Info.Durability < zeroDur {
		return
	}
	if r.Durability[core.Hydro] < zeroDur {
		return
	}
	rd := r.reduce(core.Hydro, a.Info.Durability, 0.5)
	atkDur := calcSwirlAtkDurability(rd, a.Info.Durability)
	a.Info.Durability -= rd
	//queue an attack first
	r.core.Events.Emit(core.OnSwirlHydro, r.self, a)
	r.queueSwirl(
		core.SwirlHydro,
		core.Hydro,
		core.AttackTagSwirlHydro,
		core.ICDTagSwirlHydro,
		atkDur,
		a.Info.ActorIndex,
	)
}

func (r *Reactable) trySwirlCryo(a *core.AttackEvent) {
	if a.Info.Durability < zeroDur {
		return
	}
	if r.Durability[core.Cryo] < zeroDur {
		return
	}
	rd := r.reduce(core.Cryo, a.Info.Durability, 0.5)
	atkDur := calcSwirlAtkDurability(rd, a.Info.Durability)
	a.Info.Durability -= rd
	//queue an attack first
	r.core.Events.Emit(core.OnSwirlCryo, r.self, a)
	r.queueSwirl(
		core.SwirlCryo,
		core.Cryo,
		core.AttackTagSwirlCryo,
		core.ICDTagSwirlCryo,
		atkDur,
		a.Info.ActorIndex,
	)
}

func (r *Reactable) trySwirlPyro(a *core.AttackEvent) {
	if a.Info.Durability < zeroDur {
		return
	}
	if r.Durability[core.Pyro] < zeroDur {
		return
	}
	rd := r.reduce(core.Pyro, a.Info.Durability, 0.5)
	atkDur := calcSwirlAtkDurability(rd, a.Info.Durability)
	a.Info.Durability -= rd
	//queue an attack first
	r.core.Events.Emit(core.OnSwirlPyro, r.self, a)
	r.queueSwirl(
		core.SwirlPyro,
		core.Pyro,
		core.AttackTagSwirlPyro,
		core.ICDTagSwirlPyro,
		atkDur,
		a.Info.ActorIndex,
	)
}

func (r *Reactable) trySwirlFrozen(a *core.AttackEvent) {
	if a.Info.Durability < zeroDur {
		return
	}
	if r.Durability[core.Frozen] < zeroDur {
		return
	}
	rd := r.reduce(core.Frozen, a.Info.Durability, 0.5)
	atkDur := calcSwirlAtkDurability(rd, a.Info.Durability)
	a.Info.Durability -= rd
	//queue an attack first
	r.core.Events.Emit(core.OnSwirlCryo, r.self, a)
	r.queueSwirl(
		core.SwirlCryo,
		core.Cryo,
		core.AttackTagSwirlCryo,
		core.ICDTagSwirlCryo,
		atkDur,
		a.Info.ActorIndex,
	)
	r.checkFreeze()
}
