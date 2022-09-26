package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/gadget"
)

func (r *Reactable) tryBloom(a *combat.AttackEvent) {
	//can be hydro bloom, dendro bloom, or quicken bloom
	if a.Info.Durability < ZeroDur {
		return
	}
	var consumed combat.Durability
	switch a.Info.Element {
	case attributes.Hydro:
		//this part is annoying. bloom will happen if any of the dendro like aura is present
		//so we gotta check for all 3...
		switch {
		case r.Durability[ModifierDendro] > ZeroDur:
		case r.Durability[ModifierQuicken] > ZeroDur:
		case r.Durability[ModifierBurningFuel] > ZeroDur:
		default:
			return
		}
		//reduce only check for one element so have to call twice to check for quicken as well
		consumed = r.reduce(attributes.Dendro, a.Info.Durability, 0.5)
		f := r.reduce(attributes.Quicken, a.Info.Durability, 0.5)
		if f > consumed {
			consumed = f
		}
	case attributes.Dendro:
		if r.Durability[ModifierHydro] < ZeroDur {
			return
		}
		consumed = r.reduce(attributes.Hydro, a.Info.Durability, 2)
		// case attributes.Quicken:
		// 	//TODO: ?? how to handle this??
		// 	if r.Durability[ModifierHydro] < ZeroDur {
		// 		return
		// 	}
		// 	consumed = r.reduce(attributes.Quicken, a.Info.Durability, 2)
	}
	a.Info.Durability -= consumed
	a.Info.Durability = max(a.Info.Durability, 0)
	a.Reacted = true

	// should add delay for spawning time maybe?
	r.addBloomGadget(a)
	r.core.Events.Emit(event.OnBloom, r.self, a)

	// check if quicken just added, then quicken gonna self-react with hydro if there's any hydro left
	if r.Durability[ModifierQuicken] >= ZeroDur && r.Durability[ModifierHydro] >= ZeroDur {
		hydroConsumed := r.reduce(attributes.Quicken, r.Durability[ModifierHydro], 0.5)
		r.Durability[ModifierHydro] -= hydroConsumed
		r.Durability[ModifierHydro] = max(r.Durability[ModifierHydro], 0)

		// same as above
		r.addBloomGadget(a)
		r.core.Events.Emit(event.OnBloom, r.self, a)
	}
}

type dendroCore struct {
	*gadget.Gadget
	reactable *Reactable
}

func (r *Reactable) addBloomGadget(a *combat.AttackEvent) {
	s := r.newDendroCore(a)
	r.core.Combat.AddGadget(s)
	// r.core.Combat.Log.NewEvent("bloom created", glog.LogElementEvent, a.Info.ActorIndex)
}

func (r *Reactable) newDendroCore(a *combat.AttackEvent) *dendroCore {
	s := &dendroCore{
		reactable: r,
	}

	x, y := r.self.Pos()
	// for simplicity, seeds spawn randomly within 1 radius of target
	x = x + r.core.Rand.Float64()
	y = y + r.core.Rand.Float64()
	s.Gadget = gadget.New(r.core, core.Coord{X: x, Y: y, R: 0.2})
	s.Gadget.Duration = 300
	s.Gadget.OnRemoved = func() {
		atk := combat.AttackInfo{
			ActorIndex:       a.Info.ActorIndex,
			DamageSrc:        r.self.Key(),
			Abil:             string(combat.Bloom),
			AttackTag:        combat.AttackTagBloom,
			ICDTag:           combat.ICDTagBloomDamage,
			ICDGroup:         combat.ICDGroupReactionA,
			Element:          attributes.Dendro,
			Durability:       0,
			IgnoreDefPercent: 1,
		}
		em := r.core.Player.ByIndex(a.Info.ActorIndex).Stat(attributes.EM)
		atk.FlatDmg = 2.0 * r.calcReactionDmg(atk, em)

		// self harm not active for now
		r.core.QueueAttack(atk, combat.NewCircleHit(s, 5, false, combat.TargettableEnemy), -1, 1)
	}

	return s
}

func (s *dendroCore) Tick() {
	//this is needed since gadget tick
	s.Gadget.Tick()
}

func (s *dendroCore) Type() combat.TargettableType { return combat.TargettableGadget }

func (s *dendroCore) Attack(atk *combat.AttackEvent, evt glog.Event) (float64, bool) {
	if atk.Info.Durability < ZeroDur {
		return 0, false
	}

	// only contact with pyro/electro to trigger burgeon/hyperbloom accordingly
	switch atk.Info.Element {
	case attributes.Electro:
		// trigger hyperbloom targets the nearest enemy
		// it can also do damage to player in small aoe
		ai := s.AIBloomReactionDamage(atk.Info.ActorIndex, combat.AttackTagHyperbloom, combat.ICDTagHyperbloomDamage,
			combat.ICDGroupReactionA, combat.StrikeTypeDefault, string(combat.Hyperbloom))

		// queue dmg nearest enemy
		if len(s.Core.Combat.Enemies()) <= 1 {
			s.Core.QueueAttack(ai, combat.NewDefSingleTarget(s.Core.Combat.DefaultTarget, combat.TargettableEnemy), -1, 5)
		} else {
			x, y := s.Pos()
			nearest := s.Core.Combat.EnemyByDistance(x, y, s.Index())[0]
			s.Core.QueueAttack(ai, combat.NewDefSingleTarget(nearest, combat.TargettableEnemy), -1, 5)
		}

		s.Core.Combat.RemoveGadget(s.Gadget.Index())
		s.Core.Events.Emit(event.OnHyperbloom, s, atk)
		// s.Core.Combat.Log.NewEvent("hyperbloom triggered", glog.LogElementEvent, atk.Info.ActorIndex)
	case attributes.Pyro:
		// trigger burgeon, aoe dendro damage
		// self damage
		ai := s.AIBloomReactionDamage(atk.Info.ActorIndex, combat.AttackTagBurgeon, combat.ICDTagBurgeonDamage,
			combat.ICDGroupReactionA, combat.StrikeTypeBlunt, string(combat.Burgeon))

		// self harm not active for now
		s.Core.QueueAttack(ai, combat.NewCircleHit(s, 5, false, combat.TargettableEnemy), -1, 1)
		s.Core.Combat.RemoveGadget(s.Gadget.Index())
		s.Core.Events.Emit(event.OnBurgeon, s, atk)
	default:
		return 0, false
	}

	return 0, false
}

func (s *dendroCore) ApplyDamage(*combat.AttackEvent, float64) {}

func (s *dendroCore) AIBloomReactionDamage(ActorIndex int, AttackTag combat.AttackTag,
	ICDTag combat.ICDTag, ICDGroup combat.ICDGroup, StrikeType combat.StrikeType, Abil string) combat.AttackInfo {
	ai := combat.AttackInfo{
		ActorIndex:       ActorIndex,
		DamageSrc:        s.Key(),
		Element:          attributes.Dendro,
		AttackTag:        AttackTag,
		ICDTag:           ICDTag,
		ICDGroup:         ICDGroup,
		StrikeType:       StrikeType,
		Abil:             Abil,
		IgnoreDefPercent: 1,
	}
	em := s.Core.Player.ByIndex(ActorIndex).Stat(attributes.EM)
	ai.FlatDmg = 3 * s.reactable.calcReactionDmg(ai, em)
	return ai
}
