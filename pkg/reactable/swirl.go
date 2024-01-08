package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

func calcSwirlAtkDurability(consumed, src reactions.Durability) reactions.Durability {
	if consumed < src {
		return 1.25*(0.5*consumed-1) + 25
	}
	return 1.25*(src-1) + 25
}

func (r *Reactable) queueSwirl(rt reactions.ReactionType, ele attributes.Element, tag attacks.AttackTag, icd attacks.ICDTag, dur reactions.Durability, charIndex int) {
	// swirl triggers two attacks; one self with no gauge
	// and one aoe with gauge
	ai := combat.AttackInfo{
		ActorIndex:       charIndex,
		DamageSrc:        r.self.Key(),
		Abil:             string(rt),
		AttackTag:        tag,
		ICDTag:           icd,
		ICDGroup:         attacks.ICDGroupReactionA,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          ele,
		IgnoreDefPercent: 1,
	}
	char := r.core.Player.ByIndex(charIndex)
	em := char.Stat(attributes.EM)
	flatdmg, snap := calcReactionDmg(char, ai, em)
	ai.FlatDmg = 0.6 * flatdmg
	// first attack is self no hitbox
	r.core.QueueAttackWithSnap(
		ai,
		snap,
		combat.NewSingleTargetHit(r.self.Key()),
		1,
	)
	// next is aoe - hydro swirls never do AoE damage, as they only spread the element
	if ele == attributes.Hydro {
		ai.FlatDmg = 0
	}
	ai.Durability = dur
	ai.Abil = string(rt) + " (aoe)"
	ap := combat.NewCircleHitOnTarget(r.self, nil, 5)
	ap.IgnoredKeys = []targets.TargetKey{r.self.Key()}
	r.core.QueueAttackWithSnap(
		ai,
		snap,
		ap,
		1,
	)
}

func (r *Reactable) TrySwirlElectro(a *combat.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}
	if r.Durability[Electro] < ZeroDur {
		return false
	}
	rd := r.reduce(attributes.Electro, a.Info.Durability, 0.5)
	atkDur := calcSwirlAtkDurability(rd, a.Info.Durability)
	a.Info.Durability -= rd
	a.Reacted = true
	// queue an attack first
	r.core.Events.Emit(event.OnSwirlElectro, r.self, a)

	// 0.1s gcd on swirl electro attack
	if !(r.swirlElectroGCD != -1 && r.core.F < r.swirlElectroGCD) {
		r.swirlElectroGCD = r.core.F + 0.1*60
		r.queueSwirl(
			reactions.SwirlElectro,
			attributes.Electro,
			attacks.AttackTagSwirlElectro,
			attacks.ICDTagSwirlElectro,
			atkDur,
			a.Info.ActorIndex,
		)
	}

	// at this point if any durability left, we need to check for prescence of
	// hydro in case of EC
	if a.Info.Durability > ZeroDur && r.Durability[Hydro] > ZeroDur {
		// trigger swirl hydro
		r.TrySwirlHydro(a)
		// check EC clean up
		r.checkEC()
	}
	return true
}

func (r *Reactable) TrySwirlHydro(a *combat.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}
	if r.Durability[Hydro] < ZeroDur {
		return false
	}
	rd := r.reduce(attributes.Hydro, a.Info.Durability, 0.5)
	atkDur := calcSwirlAtkDurability(rd, a.Info.Durability)
	a.Info.Durability -= rd
	a.Reacted = true
	// queue an attack first
	r.core.Events.Emit(event.OnSwirlHydro, r.self, a)

	// 0.1s gcd on swirl hydro attack
	if !(r.swirlHydroGCD != -1 && r.core.F < r.swirlHydroGCD) {
		r.swirlHydroGCD = r.core.F + 0.1*60
		r.queueSwirl(
			reactions.SwirlHydro,
			attributes.Hydro,
			attacks.AttackTagSwirlHydro,
			attacks.ICDTagSwirlHydro,
			atkDur,
			a.Info.ActorIndex,
		)
	}

	return true
}

func (r *Reactable) TrySwirlCryo(a *combat.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}
	if r.Durability[Cryo] < ZeroDur {
		return false
	}
	rd := r.reduce(attributes.Cryo, a.Info.Durability, 0.5)
	atkDur := calcSwirlAtkDurability(rd, a.Info.Durability)
	a.Info.Durability -= rd
	a.Reacted = true
	// queue an attack first
	r.core.Events.Emit(event.OnSwirlCryo, r.self, a)

	// 0.1s gcd on swirl cryo attack
	if !(r.swirlCryoGCD != -1 && r.core.F < r.swirlCryoGCD) {
		r.swirlCryoGCD = r.core.F + 0.1*60
		r.queueSwirl(
			reactions.SwirlCryo,
			attributes.Cryo,
			attacks.AttackTagSwirlCryo,
			attacks.ICDTagSwirlCryo,
			atkDur,
			a.Info.ActorIndex,
		)
	}

	return true
}

func (r *Reactable) TrySwirlPyro(a *combat.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}
	if r.Durability[Pyro] < ZeroDur {
		return false
	}
	rd := r.reduce(attributes.Pyro, a.Info.Durability, 0.5)
	atkDur := calcSwirlAtkDurability(rd, a.Info.Durability)
	a.Info.Durability -= rd
	a.Reacted = true
	r.burningCheck()
	// queue an attack first
	r.core.Events.Emit(event.OnSwirlPyro, r.self, a)

	// 0.1s gcd on swirl pyro attack
	if !(r.swirlPyroGCD != -1 && r.core.F < r.swirlPyroGCD) {
		r.swirlPyroGCD = r.core.F + 0.1*60
		r.queueSwirl(
			reactions.SwirlPyro,
			attributes.Pyro,
			attacks.AttackTagSwirlPyro,
			attacks.ICDTagSwirlPyro,
			atkDur,
			a.Info.ActorIndex,
		)
	}

	return true
}

func (r *Reactable) TrySwirlFrozen(a *combat.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}
	if r.Durability[Frozen] < ZeroDur {
		return false
	}
	rd := r.reduce(attributes.Frozen, a.Info.Durability, 0.5)
	atkDur := calcSwirlAtkDurability(rd, a.Info.Durability)
	a.Info.Durability -= rd
	a.Reacted = true
	// queue an attack first
	r.core.Events.Emit(event.OnSwirlCryo, r.self, a)

	// 0.1s gcd on swirl cryo attack
	if !(r.swirlCryoGCD != -1 && r.core.F < r.swirlCryoGCD) {
		r.swirlCryoGCD = r.core.F + 0.1*60
		r.queueSwirl(
			reactions.SwirlCryo,
			attributes.Cryo,
			attacks.AttackTagSwirlCryo,
			attacks.ICDTagSwirlCryo,
			atkDur,
			a.Info.ActorIndex,
		)
	}

	r.checkFreeze()
	return true
}
