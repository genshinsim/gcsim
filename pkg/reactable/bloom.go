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
	default:
		return
	}
	a.Info.Durability -= consumed
	a.Info.Durability = max(a.Info.Durability, 0)
	a.Reacted = true

	// TODO: re-check the frame delay
	r.core.Tasks.Add(func() {
		r.addBloomGadget(a)
	}, 45)
	r.core.Events.Emit(event.OnBloom, r.self, a)

	// if quicken just added, then quicken gonna self-react with hydro if there's any hydro left
	if r.Durability[ModifierQuicken] >= ZeroDur && r.Durability[ModifierHydro] >= ZeroDur {
		hydroConsumed := r.reduce(attributes.Quicken, r.Durability[ModifierHydro], 0.5)
		r.Durability[ModifierHydro] -= hydroConsumed
		r.Durability[ModifierHydro] = max(r.Durability[ModifierHydro], 0)

		// TODO: re-check the frame delay
		r.core.Tasks.Add(func() {
			r.addBloomGadget(a)
		}, 45)
		r.core.Events.Emit(event.OnBloom, r.self, a)
	}
}

type dendroCore struct {
	*gadget.Gadget
	reactable *Reactable // for func calcReactionDmg
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
	x = x + 2*r.core.Rand.Float64() - 1
	y = y + 2*r.core.Rand.Float64() - 1
	s.Gadget = gadget.New(r.core, core.Coord{X: x, Y: y, R: 0.2})
	s.Gadget.Duration = 300 // ??

	s.Gadget.OnRemoved = func() {
		ai := s.AIBloomReactionDamage(2, a.Info.ActorIndex, combat.AttackTagBloom, combat.ICDTagBloomDamage, combat.ICDGroupReactionA,
			combat.StrikeTypeDefault, string(combat.Bloom))
		// queue dmg
		r.core.QueueAttack(ai, combat.NewCircleHit(s.Gadget, 5, false, combat.TargettableEnemy), -1, 1)

		// queue self damage
		ai.Abil += " (self damage)"
		ai.FlatDmg = 0.05 * ai.FlatDmg
		r.core.QueueAttack(ai, combat.NewCircleHit(s.Gadget, 5, true, combat.TargettablePlayer), -1, 1)
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
		ai := s.AIBloomReactionDamage(3, atk.Info.ActorIndex, combat.AttackTagHyperbloom, combat.ICDTagHyperbloomDamage,
			combat.ICDGroupReactionA, combat.StrikeTypeDefault, string(combat.Hyperbloom))

		// queue dmg nearest enemy
		x, y := s.Gadget.Pos()
		enemies := s.Core.Combat.EnemyByDistance(x, y, s.Gadget.Key())
		if len(enemies) > 0 {
			s.Core.QueueAttack(ai, combat.NewCircleHit(s.Core.Combat.Enemy(enemies[0]), 1, false, combat.TargettableEnemy), -1, 5)

			// also queue self damage
			ai.Abil += " (self damage)"
			ai.FlatDmg = 0.05 * ai.FlatDmg
			s.Core.QueueAttack(ai, combat.NewCircleHit(s.Core.Combat.Enemy(enemies[0]), 1, true, combat.TargettablePlayer), -1, 5)
		}

		s.Core.Combat.RemoveGadget(s.Gadget.Index())
		s.Core.Events.Emit(event.OnHyperbloom, s.Gadget, atk)
		// s.Core.Combat.Log.NewEvent("hyperbloom triggered", glog.LogElementEvent, atk.Info.ActorIndex)
	case attributes.Pyro:
		// trigger burgeon, aoe dendro damage
		// self damage

		ai := s.AIBloomReactionDamage(3, atk.Info.ActorIndex, combat.AttackTagBurgeon, combat.ICDTagBurgeonDamage,
			combat.ICDGroupReactionA, combat.StrikeTypeDefault, string(combat.Burgeon))
		s.Core.QueueAttack(ai, combat.NewCircleHit(s.Gadget, 5, false, combat.TargettableEnemy), -1, 1)

		// queue self damage
		ai.Abil += " (self damage)"
		ai.FlatDmg = 0.05 * ai.FlatDmg
		s.Core.QueueAttack(ai, combat.NewCircleHit(s.Gadget, 5, true, combat.TargettablePlayer), -1, 1)

		s.Core.Combat.RemoveGadget(s.Gadget.Index())
		s.Core.Events.Emit(event.OnBurgeon, s.Gadget, atk)
	default:
		return 0, false
	}

	return 0, false
}

func (s *dendroCore) ApplyDamage(*combat.AttackEvent, float64) {}

func (s *dendroCore) AIBloomReactionDamage(BaseMultiplier float64, ActorIndex int, AttackTag combat.AttackTag,
	ICDTag combat.ICDTag, ICDGroup combat.ICDGroup, StrikeType combat.StrikeType, Abil string) combat.AttackInfo {
	ai := combat.AttackInfo{
		ActorIndex:       ActorIndex,
		DamageSrc:        s.Gadget.Key(),
		Element:          attributes.Dendro,
		AttackTag:        AttackTag,
		ICDTag:           ICDTag,
		ICDGroup:         ICDGroup,
		StrikeType:       StrikeType,
		Abil:             Abil,
		IgnoreDefPercent: 1,
	}
	em := s.Core.Player.ByIndex(ActorIndex).Stat(attributes.EM)
	ai.FlatDmg = BaseMultiplier * s.reactable.calcReactionDmg(ai, em)
	return ai
}
