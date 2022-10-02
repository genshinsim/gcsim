package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

func (r *Reactable) tryBurning(a *combat.AttackEvent) {
	if a.Info.Durability < ZeroDur {
		return
	}

	dendroDur := r.Durability[ModifierDendro]

	//adding pyro or dendro just adds to durability
	switch a.Info.Element {
	case attributes.Pyro:
		//if there's no existing pyro/burning or dendro/quicken then do nothing
		if r.Durability[ModifierDendro] < ZeroDur && r.Durability[ModifierQuicken] < ZeroDur {
			return
		}
		//add to pyro durability
		// r.attachOrRefillNormalEle(ModifierPyro, a.Info.Durability)
	case attributes.Dendro:
		//if there's no existing pyro/burning or dendro/quicken then do nothing
		if r.Durability[ModifierPyro] < ZeroDur && r.Durability[ModifierBurning] < ZeroDur {
			return
		}
		dendroDur = max(dendroDur, a.Info.Durability*0.8)
		//add to dendro durability
		// r.attachOrRefillNormalEle(ModifierDendro, a.Info.Durability)
	default:
		return
	}
	// a.Reacted = true

	if r.Durability[ModifierBurningFuel] < ZeroDur {
		r.attachBurningFuel(max(dendroDur, r.Durability[ModifierQuicken]), 1)
		r.attachBurning()

		r.core.Events.Emit(event.OnBurning, r.self, a)
		r.calcBurningDmg(a)

		if r.burningTickSrc == -1 {
			r.burningTickSrc = r.core.F
			if t, ok := r.self.(Enemy); ok {
				//queue up burning ticks
				t.QueueEnemyTask(r.nextBurningTick(r.core.F, 1, t), 15)
			}
		}
	} else {
		//overwrite burning fuel and recalc burning dmg
		if a.Info.Element == attributes.Dendro {
			r.attachBurningFuel(a.Info.Durability, 0.8)
		}
		r.calcBurningDmg(a)
	}
}

func (r *Reactable) attachBurningFuel(dur combat.Durability, mult combat.Durability) {
	//burning fuel always overwrites
	r.Durability[ModifierBurningFuel] = mult * dur
	decayRate := mult * dur / (6*dur + 420)
	if decayRate < 10.0/60.0 {
		decayRate = 10.0 / 60.0
	}
	r.DecayRate[ModifierBurningFuel] = decayRate
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
	char := r.core.Player.ByIndex(a.Info.ActorIndex)
	em := char.Stat(attributes.EM)
	atk.FlatDmg = 0.25 * calcReactionDmg(char, atk, em)
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
				combat.NewCircleHit(r.self, 1, true, r.self.Type()),
				-1,
				0,
			)
		}
		counter++
		//queue up next tick
		t.QueueEnemyTask(r.nextBurningTick(src, counter, t), 15)
	}
}

// burningCheck purges modifiers if burning no longer active
func (r *Reactable) burningCheck() {
	if r.Durability[ModifierBurning] < ZeroDur && r.Durability[ModifierBurningFuel] > ZeroDur {
		//no more burning ticks
		r.burningTickSrc = -1
		//remove burning fuel
		r.Durability[ModifierBurningFuel] = 0
		r.DecayRate[ModifierBurningFuel] = 0
	}
}
