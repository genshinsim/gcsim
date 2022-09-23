package reactable

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

func (r *Reactable) tryBurning(a *combat.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}

	//adding pyro or dendro just adds to durability
	switch a.Info.Element {
	case attributes.Pyro:
		//if there's no existing pyro/burning or dendro/quicken then do nothing
		if r.Durability[ModifierDendro] < ZeroDur && r.Durability[ModifierQuicken] < ZeroDur {
			return
		}
		//add to pyro durability
		r.attachOrRefillNormalEle(ModifierPyro, a.Info.Durability)
	case attributes.Dendro:
		//if there's no existing pyro/burning or dendro/quicken then do nothing
		if r.Durability[ModifierPyro] < ZeroDur && r.Durability[ModifierBurning] < ZeroDur {
			return
		}
		//add to dendro durability
		r.attachOrRefillNormalEle(ModifierDendro, a.Info.Durability)
	default:
		return
	}
	a.Reacted = true

	if r.Durability[ModifierBurningFuel] < ZeroDur {
		//trigger burning
		//save dendro and quicken decay rate
		if r.Durability[ModifierDendro] > ZeroDur {
			r.burningCachedDendroDecayRate = r.DecayRate[ModifierDendro]
		}
		if r.Durability[ModifierQuicken] > ZeroDur {
			r.burningCachedQuickenDecayRate = r.DecayRate[ModifierQuicken]
		}
		r.attachBurningFuel(max(r.Durability[ModifierDendro], r.Durability[ModifierQuicken]), 1)

		//update dendro and quicken decay rate to decay rate of burning fuel
		//TODO: does this also update if burning fuel gets overwritten?
		//if yes, move this to attachBurningFuel (impossible to test unless we get access to 200 durability dendro)
		if r.Durability[ModifierDendro] > ZeroDur {
			r.DecayRate[ModifierDendro] = r.DecayRate[ModifierBurningFuel]
		}
		if r.Durability[ModifierQuicken] > ZeroDur {
			r.DecayRate[ModifierDendro] = r.DecayRate[ModifierBurningFuel]
		}
		r.attachBurning()

		r.core.Events.Emit(event.OnBurning, r.self, a)
		r.calcBurningDmg(a)

		if r.burningTickSrc == -1 {
			r.burningTickSrc = r.core.F
			if t, ok := r.self.(Enemy); ok {
				//queue up burning ticks
				t.QueueEnemyTask(r.nextBurningTick(r.core.F, 1, t), 15)
				//need to reset src and restore decay rate when burning is reacted off
				r.burningEventSub()
			}
		}
	} else if r.Durability[ModifierBurningFuel] > ZeroDur {
		//overwrite burning fuel and recalc burning dmg
		if a.Info.Element == attributes.Dendro {
			r.attachBurningFuel(a.Info.Durability, 0.8)
		}
		r.calcBurningDmg(a)
	}
}

func (r *Reactable) calcBurningDmg(a *combat.AttackEvent) {
	atk := combat.AttackInfo{
		ActorIndex:       a.Info.ActorIndex,
		DamageSrc:        r.self.Index(),
		Abil:             string(combat.Burning),
		AttackTag:        combat.AttackTagBurningDamage,
		ICDTag:           combat.ICDTagBurningDamage,
		ICDGroup:         combat.ICDGroupBurning,
		Element:          attributes.Pyro,
		Durability:       25,
		IgnoreDefPercent: 1,
	}
	em := r.core.Player.ByIndex(a.Info.ActorIndex).Stat(attributes.EM)
	atk.FlatDmg = 0.25 * r.calcReactionDmg(atk, em)
	r.burningSnapshot = atk
}

func (r *Reactable) nextBurningTick(src int, counter int, t Enemy) func() {
	return func() {
		if r.burningTickSrc != src {
			//source changed, do nothing
			return
		}
		//burning SHOULD be active still, since if not we would have
		//called cleanup and set source to -1
		if r.Durability[ModifierBurningFuel] < ZeroDur || r.Durability[ModifierBurning] < ZeroDur {
			return
		}
		//so burning is active, which means both auras must still have value > 0, so we can do dmg
		if counter != 9 {
			// skip the 9th tick because hyv spaghetti
			r.core.QueueAttack(
				r.burningSnapshot,
				combat.NewCircleHit(r.self, 1, true, combat.TargettableEnemy),
				-1,
				0,
			)
		}
		counter++
		//queue up next tick
		t.QueueEnemyTask(r.nextBurningTick(src, counter, t), 15)
	}
}

func (r *Reactable) burningEventSub() {
	burningReactCheck := func(args ...interface{}) bool {
		if r.Durability[ModifierBurning] < ZeroDur && r.Durability[ModifierBurningFuel] > ZeroDur {
			//no more burning ticks
			r.burningTickSrc = -1
			//remove burning fuel
			r.Durability[ModifierBurningFuel] = 0
			r.DecayRate[ModifierBurningFuel] = 0
			//restore dendro and quicken decay rate
			if r.Durability[ModifierDendro] > ZeroDur {
				r.DecayRate[ModifierDendro] = r.burningCachedDendroDecayRate
				r.burningCachedDendroDecayRate = 0
			}
			if r.Durability[ModifierQuicken] > ZeroDur {
				r.DecayRate[ModifierQuicken] = r.burningCachedQuickenDecayRate
				r.burningCachedQuickenDecayRate = 0
			}
			// remove react check
			r.core.Events.Unsubscribe(event.OnVaporize, fmt.Sprintf("burning-vaporize-%v", r.self.Index()))
			r.core.Events.Unsubscribe(event.OnOverload, fmt.Sprintf("burning-overload-%v", r.self.Index()))
			r.core.Events.Unsubscribe(event.OnMelt, fmt.Sprintf("burning-melt-%v", r.self.Index()))
			r.core.Events.Unsubscribe(event.OnSwirlPyro, fmt.Sprintf("burning-swirlpyro-%v", r.self.Index()))
			r.core.Events.Unsubscribe(event.OnCrystallizePyro, fmt.Sprintf("burning-crystallizepyro-%v", r.self.Index()))
		}
		return false
	}
	// add react check
	r.core.Events.Subscribe(event.OnVaporize, burningReactCheck, fmt.Sprintf("burning-vaporize-%v", r.self.Index()))
	r.core.Events.Subscribe(event.OnOverload, burningReactCheck, fmt.Sprintf("burning-overload-%v", r.self.Index()))
	r.core.Events.Subscribe(event.OnMelt, burningReactCheck, fmt.Sprintf("burning-melt-%v", r.self.Index()))
	r.core.Events.Subscribe(event.OnSwirlPyro, burningReactCheck, fmt.Sprintf("burning-swirlpyro-%v", r.self.Index()))
	r.core.Events.Subscribe(event.OnCrystallizePyro, burningReactCheck, fmt.Sprintf("burning-crystallizepyro-%v", r.self.Index()))
}
