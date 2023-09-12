package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

func (r *Reactable) TryBurning(a *combat.AttackEvent) bool {
	if a.Info.Durability < ZeroDur {
		return false
	}

	dendroDur := r.Durability[Dendro]

	// adding pyro or dendro just adds to durability
	switch a.Info.Element {
	case attributes.Pyro:
		// if there's no existing pyro/burning or dendro/quicken then do nothing
		if r.Durability[Dendro] < ZeroDur && r.Durability[Quicken] < ZeroDur {
			return false
		}
		// add to pyro durability
		// r.attachOrRefillNormalEle(ModifierPyro, a.Info.Durability)
	case attributes.Dendro:
		// if there's no existing pyro/burning or dendro/quicken then do nothing
		if r.Durability[Pyro] < ZeroDur && r.Durability[Burning] < ZeroDur {
			return false
		}
		dendroDur = max(dendroDur, a.Info.Durability*0.8)
		// add to dendro durability
		// r.attachOrRefillNormalEle(ModifierDendro, a.Info.Durability)
	default:
		return false
	}
	// a.Reacted = true

	if r.Durability[BurningFuel] < ZeroDur {
		r.attachBurningFuel(max(dendroDur, r.Durability[Quicken]), 1)
		r.attachBurning()

		r.core.Events.Emit(event.OnBurning, r.self, a)
		r.calcBurningDmg(a)

		if r.burningTickSrc == -1 {
			r.burningTickSrc = r.core.F
			if t, ok := r.self.(Enemy); ok {
				// queue up burning ticks
				t.QueueEnemyTask(r.nextBurningTick(r.core.F, 1, t), 15)
			}
		}
		return true
	}
	// overwrite burning fuel and recalc burning dmg
	if a.Info.Element == attributes.Dendro {
		r.attachBurningFuel(a.Info.Durability, 0.8)
	}
	r.calcBurningDmg(a)

	return false
}

func (r *Reactable) attachBurningFuel(dur, mult reactions.Durability) {
	// burning fuel always overwrites
	r.Durability[BurningFuel] = mult * dur
	decayRate := mult * dur / (6*dur + 420)
	if decayRate < 10.0/60.0 {
		decayRate = 10.0 / 60.0
	}
	r.DecayRate[BurningFuel] = decayRate
}

func (r *Reactable) calcBurningDmg(a *combat.AttackEvent) {
	atk := combat.AttackInfo{
		ActorIndex:       a.Info.ActorIndex,
		DamageSrc:        r.self.Key(),
		Abil:             string(reactions.Burning),
		AttackTag:        attacks.AttackTagBurningDamage,
		ICDTag:           attacks.ICDTagBurningDamage,
		ICDGroup:         attacks.ICDGroupBurning,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Pyro,
		Durability:       25,
		IgnoreDefPercent: 1,
	}
	char := r.core.Player.ByIndex(a.Info.ActorIndex)
	em := char.Stat(attributes.EM)
	flatdmg, snap := calcReactionDmg(char, atk, em)
	atk.FlatDmg = 0.25 * flatdmg
	r.burningAtk = atk
	r.burningSnapshot = snap
}

func (r *Reactable) nextBurningTick(src, counter int, t Enemy) func() {
	return func() {
		if r.burningTickSrc != src {
			// source changed, do nothing
			return
		}
		// burning SHOULD be active still, since if not we would have
		// called cleanup and set source to -1
		if r.Durability[BurningFuel] < ZeroDur || r.Durability[Burning] < ZeroDur {
			return
		}
		// so burning is active, which means both auras must still have value > 0, so we can do dmg
		if counter != 9 {
			// skip the 9th tick because hyv spaghetti
			ai := r.burningAtk
			ap := combat.NewCircleHitOnTarget(r.self, nil, 1)
			r.core.QueueAttackWithSnap(
				ai,
				r.burningSnapshot,
				ap,
				0,
			)
			// self damage
			ai.Abil += " (self damage)"
			ap.SkipTargets[targets.TargettablePlayer] = false
			ap.SkipTargets[targets.TargettableEnemy] = true
			ap.SkipTargets[targets.TargettableGadget] = true
			r.core.QueueAttackWithSnap(
				ai,
				r.burningSnapshot,
				ap,
				0,
			)
		}
		counter++
		// queue up next tick
		t.QueueEnemyTask(r.nextBurningTick(src, counter, t), 15)
	}
}

// burningCheck purges modifiers if burning no longer active
func (r *Reactable) burningCheck() {
	if r.Durability[Burning] < ZeroDur && r.Durability[BurningFuel] > ZeroDur {
		// no more burning ticks
		r.burningTickSrc = -1
		// remove burning fuel
		r.Durability[BurningFuel] = 0
		r.DecayRate[BurningFuel] = 0
	}
}
