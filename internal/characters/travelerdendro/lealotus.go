package travelerdendro

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/gadget"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

type LeaLotus struct {
	*gadget.Gadget
	*reactable.Reactable
	burstAtk        *combat.AttackEvent
	char            *char
	targetingRadius float64
	hitboxRadius    float64
}

func (s *LeaLotus) AuraContains(e ...attributes.Element) bool {
	for ele := range e {
		if s.Reactable.Durability[ele] <= reactable.ZeroDur {
			return false
		}
	}
	return true
}

func (c *char) newLeaLotusLamp() *LeaLotus {
	s := &LeaLotus{}
	x, y := c.Core.Combat.Player().Pos()
	// TODO The gadget spawns 1 unit away from the player (in the direction the player is facing?)
	s.Gadget = gadget.New(c.Core, core.Coord{X: x, Y: y, R: 1}, combat.GadgetTypLeaLotus)
	s.Reactable = &reactable.Reactable{}
	s.Reactable.Init(s, c.Core)
	s.Durability[reactable.ModifierDendro] = 10

	s.Duration = 12 * 60
	if c.Base.Cons >= 2 {
		s.Duration += 3 * 60
	}

	//burst status last the duration of the gadget but is removed if pyro applied
	c.Core.Status.Add(burstKey, s.Duration)

	// First hitmark is 37f after spawn, all other pre-transfig hits will be 90f between.
	c.Core.Tasks.Add(func() {
		if !s.Alive {
			return
		}
		s.QueueAttack(0)
		//repeat attack every 90
		s.Gadget.OnThinkInterval = func() {
			s.QueueAttack(0)
		}
		s.Gadget.ThinkInterval = 90
	}, burstHitmark-leaLotusAppear)

	c.burstTransfig = attributes.NoElement
	s.char = c

	s.targetingRadius = 8
	s.hitboxRadius = 2
	c.burstOverflowingLotuslight = 0

	procAI := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lea Lotus Lamp",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       burstTick[c.TalentLvlBurst()],
	}
	s.burstAtk = &combat.AttackEvent{
		Info:     procAI,
		Snapshot: c.Snapshot(&procAI),
	}

	return s
}

func (s *LeaLotus) HandleAttack(atk *combat.AttackEvent) float64 {
	s.Core.Events.Emit(event.OnGadgetHit, s, atk)

	s.Core.Log.NewEvent(fmt.Sprintf("dmc lamp hit by %s", atk.Info.Abil), glog.LogCharacterEvent, s.char.Index)

	s.ShatterCheck(atk)

	if atk.Info.Durability > 0 {
		atk.Info.Durability *= combat.Durability(s.WillApplyEle(atk.Info.ICDTag, atk.Info.ICDGroup, atk.Info.ActorIndex))
		if atk.Info.Durability > 0 && atk.Info.Element != attributes.Physical {
			existing := s.Reactable.ActiveAuraString()
			applied := atk.Info.Durability
			s.React(atk)
			if s.Core.Flags.LogDebug && atk.Reacted {
				s.Core.Log.NewEvent(
					"application",
					glog.LogElementEvent,
					atk.Info.ActorIndex,
				).
					Write("attack_tag", atk.Info.AttackTag).
					Write("applied_ele", atk.Info.Element.String()).
					Write("dur", applied).
					Write("abil", atk.Info.Abil).
					Write("target", s.Key()).
					Write("existing", existing).
					Write("after", s.Reactable.ActiveAuraString())

			}
		}
	}

	//apply damage delay is only there to make sure aura gets applied at the end of current frame
	//however because we can only hold cryo, we'll only call this if atk is cryo and there
	//is durability left
	if atk.Info.Element != attributes.Cryo {
		return 0
	}
	if atk.Info.Durability < reactable.ZeroDur {
		return 0
	}
	if atk.Reacted {
		return 0
	}
	s.Core.Combat.Tasks.Add(func() {
		s.attachEle(atk)
	}, 0)
	return 0
}

func (s *LeaLotus) attachEle(atk *combat.AttackEvent) {
	// check for ICD first
	existing := s.Reactable.ActiveAuraString()
	applied := atk.Info.Durability
	s.AttachOrRefill(atk)
	if s.Core.Flags.LogDebug {
		s.Core.Log.NewEvent(
			"application",
			glog.LogElementEvent,
			atk.Info.ActorIndex,
		).
			Write("attack_tag", atk.Info.AttackTag).
			Write("applied_ele", atk.Info.Element.String()).
			Write("dur", applied).
			Write("abil", atk.Info.Abil).
			Write("target", s.Key()).
			Write("existing", existing).
			Write("after", s.Reactable.ActiveAuraString())

	}
}

func (s *LeaLotus) Tick() {
	//this is needed since gadget tick
	s.Reactable.Tick()
	s.Gadget.Tick()
}

func (l *LeaLotus) QueueAttack(delay int) {
	x, y := l.Gadget.Pos()
	enemies := l.Core.Combat.EnemiesWithinRadius(x, y, l.targetingRadius)
	if len(enemies) > 0 {
		idx := l.Core.Rand.Intn(len(enemies))

		l.Core.QueueAttackWithSnap(
			l.burstAtk.Info,
			l.burstAtk.Snapshot,
			combat.NewCircleHit(l.Core.Combat.Enemy(enemies[idx]), l.hitboxRadius),
			delay,
		)
	}
}

func (r *LeaLotus) React(a *combat.AttackEvent) {
	//only check the ones possible
	switch a.Info.Element {
	case attributes.Electro:
		r.TryAggravate(a)
		r.TryFrozenSuperconduct(a)
		r.TrySuperconduct(a)
		r.TryQuicken(a)
	case attributes.Pyro:
		r.TryMelt(a)
		r.TryBurning(a)
	case attributes.Cryo:
	case attributes.Hydro:
		r.TryFreeze(a)
		r.TryBloom(a)
	case attributes.Anemo:
		r.TrySwirlHydro(a)
		r.TrySwirlCryo(a)
		r.TrySwirlFrozen(a)
	case attributes.Geo:
		r.TryCrystallizeCryo(a)
		r.TryCrystallizeFrozen(a)
	case attributes.Dendro:
		r.TrySpread(a)
		r.TryBloom(a)
	}
}

func (s *LeaLotus) TryQuicken(a *combat.AttackEvent) {
	if !s.Reactable.TryQuicken(a) {
		return
	}
	for t := 15; t <= s.Duration; t += 54 {
		s.QueueAttack(t)
	}
	s.transfig(attributes.Electro)
}

func (s *LeaLotus) TryBloom(a *combat.AttackEvent) {
	if !s.Reactable.TryBloom(a) {
		return
	}
	s.targetingRadius = 12
	s.hitboxRadius = 4
	for t := 15; t <= s.Duration; t += 90 {
		s.QueueAttack(t)
	}
	s.transfig(attributes.Hydro)
}

func (s *LeaLotus) TryBurning(a *combat.AttackEvent) {
	if !s.Reactable.TryBurning(a) {
		return
	}
	s.burstAtk.Info.Abil = "Lea Lotus Lamp Explosion"
	s.burstAtk.Info.Durability = 50
	s.burstAtk.Info.ICDTag = combat.ICDTagNone
	s.burstAtk.Info.Mult = burstExplode[s.char.TalentLvlBurst()]
	s.Core.Tasks.Add(func() {
		s.Core.QueueAttackWithSnap(s.burstAtk.Info, s.burstAtk.Snapshot, combat.NewCircleHit(s, 6.5), 0)
		s.Core.Status.Delete(burstKey)
	}, 60)
	s.transfig(attributes.Pyro)
}

func (s *LeaLotus) transfig(ele attributes.Element) {
	s.Core.Log.NewEvent(fmt.Sprintf("dmc lamp %s transfig triggered", ele.String()), glog.LogCharacterEvent, s.char.Index)
	s.char.burstTransfig = ele
	if s.char.Base.Cons >= 4 {
		s.char.c4()
	}
	s.Kill()
}
