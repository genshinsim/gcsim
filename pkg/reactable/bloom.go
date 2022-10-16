package reactable

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/gadget"
)

const DendroCoreDelay = 20

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

	r.addBloomGadget(a)
	r.core.Events.Emit(event.OnBloom, r.self, a)
}

// this function should only be called after a catalyze reaction (queued to the end of current frame)
// this reaction will check if any hydro exists and if so trigger a bloom reaction
func (r *Reactable) tryQuickenBloom(a *combat.AttackEvent) {
	if r.Durability[ModifierQuicken] < ZeroDur {
		//this should be a sanity check; should not happen realistically unless something wipes off
		//the quicken immediately (same frame) after catalyze
		return
	}
	if r.Durability[ModifierHydro] < ZeroDur {
		return
	}
	avail := r.Durability[ModifierQuicken]
	consumed := r.reduce(attributes.Hydro, avail, 2)
	r.Durability[ModifierQuicken] -= consumed

	r.addBloomGadget(a)
	r.core.Events.Emit(event.OnBloom, r.self, a)
}

type DendroCore struct {
	*gadget.Gadget
	srcFrame int
}

func (r *Reactable) addBloomGadget(a *combat.AttackEvent) {
	r.core.Tasks.Add(func() {
		var t combat.Gadget = NewDendroCore(r.core, r.self, a)
		r.core.Combat.AddGadget(t)
		r.core.Events.Emit(event.OnDendroCore, t, a)
	}, DendroCoreDelay)
}

func NewDendroCore(c *core.Core, pos combat.Positional, a *combat.AttackEvent) *DendroCore {
	s := &DendroCore{
		srcFrame: c.F,
	}

	x, y := pos.Pos()
	// for simplicity, seeds spawn randomly within 1 radius of target
	x = x + 2*c.Rand.Float64() - 1
	y = y + 2*c.Rand.Float64() - 1
	s.Gadget = gadget.New(c, core.Coord{X: x, Y: y, R: 0.2}, combat.GadgetTypDendroCore)
	s.Gadget.Duration = 300 // ??

	char := s.Core.Player.ByIndex(a.Info.ActorIndex)

	explode := func() {
		ai := NewBloomAttack(char, s)
		c.QueueAttack(ai, combat.NewCircleHit(s, 5, false, combat.TargettableEnemy), -1, 1)

		//self damage
		ai.Abil += " (self damage)"
		ai.FlatDmg = 0.05 * ai.FlatDmg
		c.QueueAttack(ai, combat.NewCircleHit(s.Gadget, 5, true, combat.TargettablePlayer), -1, 1)
	}
	//TODO: should bloom do damage if it blows up due to limit reached?
	s.Gadget.OnExpiry = explode
	s.Gadget.OnKill = explode

	return s
}

func (s *DendroCore) Tick() {
	//this is needed since gadget tick
	s.Gadget.Tick()
}

func (s *DendroCore) Attack(atk *combat.AttackEvent, evt glog.Event) (float64, bool) {
	if atk.Info.Durability < ZeroDur {
		return 0, false
	}

	char := s.Core.Player.ByIndex(atk.Info.ActorIndex)
	// only contact with pyro/electro to trigger burgeon/hyperbloom accordingly
	switch atk.Info.Element {
	case attributes.Electro:
		// trigger hyperbloom targets the nearest enemy
		// it can also do damage to player in small aoe
		ai := NewHyperbloomAttack(char, s)
		// queue dmg nearest enemy
		x, y := s.Gadget.Pos()
		enemies := s.Core.Combat.EnemyByDistance(x, y, combat.InvalidTargetKey)
		if len(enemies) > 0 {
			s.Core.QueueAttack(ai, combat.NewCircleHit(s.Core.Combat.Enemy(enemies[0]), 1, false, combat.TargettableEnemy), -1, 5)

			// also queue self damage
			ai.Abil += " (self damage)"
			ai.FlatDmg = 0.05 * ai.FlatDmg
			s.Core.QueueAttack(ai, combat.NewCircleHit(s.Core.Combat.Enemy(enemies[0]), 1, true, combat.TargettablePlayer), -1, 5)
		}

		s.Gadget.OnKill = nil
		s.Gadget.Kill()
		s.Core.Events.Emit(event.OnHyperbloom, s, atk)
	case attributes.Pyro:
		// trigger burgeon, aoe dendro damage
		// self damage
		ai := NewBurgeonAttack(char, s)

		s.Core.QueueAttack(ai, combat.NewCircleHit(s.Gadget, 5, false, combat.TargettableEnemy), -1, 1)

		// queue self damage
		ai.Abil += " (self damage)"
		ai.FlatDmg = 0.05 * ai.FlatDmg
		s.Core.QueueAttack(ai, combat.NewCircleHit(s.Gadget, 5, true, combat.TargettablePlayer), -1, 1)

		s.Gadget.OnKill = nil
		s.Gadget.Kill()
		s.Core.Events.Emit(event.OnBurgeon, s, atk)
	default:
		return 0, false
	}

	return 0, false
}

func (s *DendroCore) ApplyDamage(*combat.AttackEvent, float64) {}

const (
	BloomMultiplier      = 2
	BurgeonMultiplier    = 3
	HyperbloomMultiplier = 3
)

func NewBloomAttack(char *character.CharWrapper, src combat.Target) combat.AttackInfo {
	em := char.Stat(attributes.EM)
	ai := combat.AttackInfo{
		ActorIndex:       char.Index,
		DamageSrc:        src.Key(),
		Element:          attributes.Dendro,
		AttackTag:        combat.AttackTagBloom,
		ICDTag:           combat.ICDTagBloomDamage,
		ICDGroup:         combat.ICDGroupReactionA,
		StrikeType:       combat.StrikeTypeDefault,
		Abil:             string(combat.Bloom),
		IgnoreDefPercent: 1,
	}
	ai.FlatDmg = BloomMultiplier * calcReactionDmg(char, ai, em)
	return ai
}

func NewBurgeonAttack(char *character.CharWrapper, src combat.Target) combat.AttackInfo {
	em := char.Stat(attributes.EM)
	ai := combat.AttackInfo{
		ActorIndex:       char.Index,
		DamageSrc:        src.Key(),
		Element:          attributes.Dendro,
		AttackTag:        combat.AttackTagBurgeon,
		ICDTag:           combat.ICDTagBurgeonDamage,
		ICDGroup:         combat.ICDGroupReactionA,
		StrikeType:       combat.StrikeTypeDefault,
		Abil:             string(combat.Burgeon),
		IgnoreDefPercent: 1,
	}
	ai.FlatDmg = BurgeonMultiplier * calcReactionDmg(char, ai, em)
	return ai
}

func NewHyperbloomAttack(char *character.CharWrapper, src combat.Target) combat.AttackInfo {
	em := char.Stat(attributes.EM)
	ai := combat.AttackInfo{
		ActorIndex:       char.Index,
		DamageSrc:        src.Key(),
		Element:          attributes.Dendro,
		AttackTag:        combat.AttackTagHyperbloom,
		ICDTag:           combat.ICDTagHyperbloomDamage,
		ICDGroup:         combat.ICDGroupReactionA,
		StrikeType:       combat.StrikeTypeDefault,
		Abil:             string(combat.Hyperbloom),
		IgnoreDefPercent: 1,
	}
	ai.FlatDmg = HyperbloomMultiplier * calcReactionDmg(char, ai, em)
	return ai
}
