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
	burstAtk *combat.AttackEvent
	char     *char
	r        reactable.Reactable
}

func (s *LeaLotus) transfigInit() {
	s.Core.Events.Subscribe(event.OnQuicken, func(args ...interface{}) bool {
		target := args[0].(combat.Target)
		if target == s {
			s.transfig(attributes.Electro)
			return true
		}
		return false
	}, "lealotus-electro")

	s.Core.Events.Subscribe(event.OnBloom, func(args ...interface{}) bool {
		target := args[0].(combat.Target)
		if target == s {
			s.transfig(attributes.Hydro)
			return true
		}
		return false
	}, "lealotus-hydro")

	s.Core.Events.Subscribe(event.OnBurning, func(args ...interface{}) bool {
		target := args[0].(combat.Target)
		if target == s {
			s.transfig(attributes.Pyro)
			return true
		}
		return false
	}, "lealotus-pyro")
}

func (s *LeaLotus) AuraContains(e ...attributes.Element) bool {
	for ele := range e {
		if s.Reactable.Durability[ele] <= reactable.ZeroDur {
			return false
		}
	}
	return true
}

func (c *char) newLeaLotusLamp(duration int) *LeaLotus {
	s := &LeaLotus{}
	x, y := c.Core.Combat.Player().Pos()
	s.Gadget = gadget.New(c.Core, core.Coord{X: x, Y: y, R: 1}, combat.GadgetTypLeaLotus)
	s.Duration = duration
	s.ThinkInterval = 90
	s.Gadget.OnThinkInterval = s.OnThinkInterval

	s.Reactable = &reactable.Reactable{}
	s.Reactable.Init(s, c.Core)
	s.Durability[reactable.ModifierDendro] = 20

	s.char = c
	s.char.burstTransfig = attributes.NoElement

	procAI := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lea Lotus Lamp (Q)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       burstTick[c.TalentLvlBurst()],
	}
	burstSnap := c.Snapshot(&procAI)
	s.burstAtk = &combat.AttackEvent{
		Info:     procAI,
		Snapshot: burstSnap,
	}
	s.char.burstAlive = true
	s.transfigInit()
	return s
}

func (s *LeaLotus) transfig(ele attributes.Element) {
	s.Core.Log.NewEvent(fmt.Sprintf("DMC lamp %s transfig triggered", ele.String()), glog.LogCharacterEvent, s.char.Index)
	s.char.burstTransfig = ele
	switch ele {
	case attributes.Hydro:
		for t := 15; t <= s.Duration; t += 90 {
			s.Core.QueueAttackWithSnap(s.burstAtk.Info, s.burstAtk.Snapshot, combat.NewCircleHit(s.Core.Combat.PrimaryTarget(), 5, false, combat.TargettableEnemy), t)
		}
		s.Core.Tasks.Add(func() {
			s.char.burstAlive = false
		}, s.Gadget.Duration)

	case attributes.Electro:
		for t := 15; t <= s.Duration; t += 54 {
			s.Core.QueueAttackWithSnap(s.burstAtk.Info, s.burstAtk.Snapshot, combat.NewCircleHit(s.Core.Combat.PrimaryTarget(), 3, false, combat.TargettableEnemy), t)
		}
		s.Core.Tasks.Add(func() {
			s.char.burstAlive = false
		}, s.Gadget.Duration)

	case attributes.Pyro:
		s.burstAtk.Info.Abil = "Lea Lotus Lamp Explosion (Q)"
		s.burstAtk.Info.Durability = 50
		s.burstAtk.Info.ICDTag = combat.ICDTagNone
		s.burstAtk.Info.Mult = burstExplode[s.char.TalentLvlBurst()]
		s.Core.QueueAttackWithSnap(s.burstAtk.Info, s.burstAtk.Snapshot, combat.NewCircleHit(s.Core.Combat.PrimaryTarget(), 5, false, combat.TargettableEnemy), 60)
		s.Core.Tasks.Add(func() {
			s.char.burstAlive = false
			s.Core.Status.Delete("dmc-burst") // starts on first hitmark
		}, 60)
	}
	if s.char.Base.Cons >= 4 {
		s.char.c4()
	}
	s.Kill()
}

func (s *LeaLotus) Attack(a *combat.AttackEvent, evt glog.Event) (float64, bool) {
	s.React(a)
	return 0, false
}

func (s *LeaLotus) Tick() {
	//this is needed since gadget tick
	s.Reactable.Tick()
	s.Gadget.Tick()
}

func (s *LeaLotus) OnThinkInterval() {
	s.Core.QueueAttackWithSnap(s.burstAtk.Info, s.burstAtk.Snapshot, combat.NewCircleHit(s.Core.Combat.PrimaryTarget(), 3, false, combat.TargettableEnemy), 0)
}

func (s *LeaLotus) ApplyDamage(*combat.AttackEvent, float64) {}

func (s *LeaLotus) OnExpiry(*combat.AttackEvent, float64) {
	s.char.burstAlive = false
	// s.Core.Events.Unsubscribe(event.OnQuicken, "lealotus-electro")
	// s.Core.Events.Unsubscribe(event.OnBloom, "lealotus-hydro")
	// s.Core.Events.Unsubscribe(event.OnBurning, "lealotus-pyro")
}

func (s *LeaLotus) OnKill(*combat.AttackEvent, float64) {
	// s.Core.Events.Unsubscribe(event.OnQuicken, "lealotus-electro")
	// s.Core.Events.Unsubscribe(event.OnBloom, "lealotus-hydro")
	// s.Core.Events.Unsubscribe(event.OnBurning, "lealotus-pyro")
}
