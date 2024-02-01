package dendro

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
	"github.com/genshinsim/gcsim/pkg/gadget"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

type LeaLotus struct {
	*gadget.Gadget
	*reactable.Reactable
	burstAtk     *combat.AttackEvent
	char         *Traveler
	hitboxRadius float64
}

func (c *Traveler) newLeaLotusLamp() *LeaLotus {
	s := &LeaLotus{}
	player := c.Core.Combat.Player()
	c.burstPos = geometry.CalcOffsetPoint(
		player.Pos(),
		geometry.Point{Y: 1},
		player.Direction(),
	)
	s.Gadget = gadget.New(c.Core, c.burstPos, 1, combat.GadgetTypLeaLotus)
	s.Reactable = &reactable.Reactable{}
	s.Reactable.Init(s, c.Core)
	s.Durability[reactable.Dendro] = 10

	s.Duration = 12 * 60
	if c.Base.Cons >= 2 {
		s.Duration += 3 * 60
	}

	// burst status last the duration of the gadget but is removed if pyro applied
	c.Core.Status.Add(burstKey, s.Duration)

	// First hitmark is 37f after spawn, all other pre-transfig hits will be 90f between.
	c.Core.Tasks.Add(func() {
		if !s.Alive {
			return
		}
		s.QueueAttack(0)
		// repeat attack every 90
		s.Gadget.OnThinkInterval = func() {
			s.QueueAttack(0)
		}
		s.Gadget.ThinkInterval = 90
	}, burstHitmark-leaLotusAppear)

	c.burstTransfig = attributes.NoElement
	s.char = c

	c.burstRadius = 8
	s.hitboxRadius = 2
	c.burstOverflowingLotuslight = 0

	procAI := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lea Lotus Lamp",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
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

	s.PoiseDMGCheck(atk)
	s.ShatterCheck(atk)

	if atk.Info.Durability > 0 {
		atk.Info.Durability *= reactions.Durability(s.WillApplyEle(atk.Info.ICDTag, atk.Info.ICDGroup, atk.Info.ActorIndex))
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

	// apply damage delay is only there to make sure aura gets applied at the end of current frame
	// however because we can only hold cryo, we'll only call this if atk is cryo and there
	// is durability left
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
	// this is needed since gadget tick
	s.Reactable.Tick()
	s.Gadget.Tick()
}

func (s *LeaLotus) QueueAttack(delay int) {
	enemy := s.Core.Combat.RandomEnemyWithinArea(combat.NewCircleHitOnTarget(s.Gadget, nil, s.char.burstRadius), nil)
	if enemy == nil {
		return
	}
	s.Core.QueueAttackWithSnap(
		s.burstAtk.Info,
		s.burstAtk.Snapshot,
		combat.NewCircleHitOnTarget(enemy, nil, s.hitboxRadius),
		delay,
	)
}

func (s *LeaLotus) React(a *combat.AttackEvent) {
	// only check the ones possible
	switch a.Info.Element {
	case attributes.Electro:
		s.TryAggravate(a)
		s.TryFrozenSuperconduct(a)
		s.TrySuperconduct(a)
		s.TryQuicken(a)
	case attributes.Pyro:
		s.TryMelt(a)
		s.TryBurning(a)
	case attributes.Cryo:
	case attributes.Hydro:
		s.TryFreeze(a)
		s.TryBloom(a)
	case attributes.Anemo:
		s.TrySwirlHydro(a)
		s.TrySwirlCryo(a)
		s.TrySwirlFrozen(a)
	case attributes.Geo:
		s.TryCrystallizeCryo(a)
		s.TryCrystallizeFrozen(a)
	case attributes.Dendro:
		s.TrySpread(a)
		s.TryBloom(a)
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
	s.char.burstRadius = 12
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
	s.burstAtk.Info.ICDTag = attacks.ICDTagNone
	s.burstAtk.Info.Mult = burstExplode[s.char.TalentLvlBurst()]
	s.Core.Tasks.Add(func() {
		s.Core.QueueAttackWithSnap(
			s.burstAtk.Info,
			s.burstAtk.Snapshot,
			combat.NewCircleHitOnTarget(s, nil, 6.5),
			0,
		)
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

func (s *LeaLotus) SetDirection(trg geometry.Point) {}
func (s *LeaLotus) SetDirectionToClosestEnemy()     {}
func (s *LeaLotus) CalcTempDirection(trg geometry.Point) geometry.Point {
	return geometry.DefaultDirection()
}
