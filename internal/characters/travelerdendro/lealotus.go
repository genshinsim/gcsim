package travelerdendro

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/enemy"
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
	r               reactable.Reactable
	//icd related
	icdTagOnTimer       [enemy.MaxTeamSize][combat.ICDTagLength]bool
	icdTagCounter       [enemy.MaxTeamSize][combat.ICDTagLength]int
	icdDamageTagOnTimer [enemy.MaxTeamSize][combat.ICDTagLength]bool
	icdDamageTagCounter [enemy.MaxTeamSize][combat.ICDTagLength]int
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
	s.ThinkInterval = 37
	s.Gadget.OnThinkInterval = s.OnThinkInterval

	s.Reactable = &reactable.Reactable{}
	s.Reactable.Init(s, c.Core)
	s.Durability[reactable.ModifierDendro] = 20

	s.char = c
	s.char.burstTransfig = attributes.NoElement

	s.targetingRadius = 8
	s.hitboxRadius = 3

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
		s.targetingRadius = 12
		s.hitboxRadius = 5
		for t := 15; t <= s.Duration; t += 90 {
			s.QueueAttack(t)
		}
		s.Core.Tasks.Add(func() {
			s.char.burstAlive = false
		}, s.Gadget.Duration)

	case attributes.Electro:
		for t := 15; t <= s.Duration; t += 54 {
			s.QueueAttack(t)
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

	s.Core.Log.NewEvent(fmt.Sprintf("DMC lamp hit by %s", a.Info.Abil), glog.LogCharacterEvent, s.char.Index)

	s.ShatterCheck(a)
	//check for ICD first
	a.OnICD = !s.WillApplyEle(a.Info.ICDTag, a.Info.ICDGroup, a.Info.ActorIndex)
	if a.Info.Durability > 0 && !a.OnICD && a.Info.Element != attributes.Physical {
		s.React(a)
	}
	return 0, false
}

func (s *LeaLotus) Tick() {
	//this is needed since gadget tick
	s.Reactable.Tick()
	s.Gadget.Tick()
}

func (s *LeaLotus) OnThinkInterval() {
	s.ThinkInterval = 90
	s.QueueAttack(0)
}

func (s *LeaLotus) ApplyDamage(atk *combat.AttackEvent, damage float64) {
	if atk.Info.Durability > 0 && !atk.OnICD && atk.Info.Element != attributes.Physical {
		if !atk.Reacted {
			s.AttachOrRefill(atk)
			s.Core.Log.NewEvent(fmt.Sprintf("Aura Applied: DMC lamp auras: %s", s.ActiveAuraString()), glog.LogCharacterEvent, s.char.Index)
		}

		s.Core.Log.NewEvent(fmt.Sprintf("DMC lamp auras: %s", s.ActiveAuraString()), glog.LogCharacterEvent, s.char.Index)

	}
}

func (s *LeaLotus) OnExpiry(*combat.AttackEvent, float64) {
	s.char.burstAlive = false
}

func (l *LeaLotus) QueueAttack(delay int) {
	x, y := l.Gadget.Pos()
	enemies := l.Core.Combat.EnemiesWithinRadius(x, y, l.targetingRadius)
	if len(enemies) > 0 {
		idx := l.Core.Rand.Intn(len(enemies))

		l.Core.QueueAttackWithSnap(
			l.burstAtk.Info,
			l.burstAtk.Snapshot,
			combat.NewCircleHit(l.Core.Combat.Enemy(enemies[idx]), l.hitboxRadius, false, combat.TargettableEnemy),
			delay,
		)
	}
}

// ICD code copied from enemy/icd.go
func (t *LeaLotus) WillApplyEle(tag combat.ICDTag, grp combat.ICDGroup, char int) bool {

	//no icd if no tag
	if tag == combat.ICDTagNone {
		return true
	}

	//check if we need to start timer
	x := t.icdTagOnTimer[char][tag]
	if !t.icdTagOnTimer[char][tag] {
		t.icdTagOnTimer[char][tag] = true
		t.ResetTagCounterAfterDelay(tag, grp, char)
	}

	val := t.icdTagCounter[char][tag]
	//increment the counter
	//if counter > length, then reset back to 0
	t.icdTagCounter[char][tag]++
	if t.icdTagCounter[char][tag] == len(combat.ICDGroupEleApplicationSequence[grp]) {
		t.icdTagCounter[char][tag] = 0
	}

	t.Core.Log.NewEvent("ele icd check", glog.LogICDEvent, char).
		Write("grp", grp).
		Write("target", t.TargetIndex).
		Write("tag", tag).
		Write("counter", val).
		Write("val", combat.ICDGroupEleApplicationSequence[grp][val]).
		Write("group on timer", x)
	//true if group seq is 1
	return combat.ICDGroupEleApplicationSequence[grp][val] == 1
}

func (t *LeaLotus) GroupTagDamageMult(tag combat.ICDTag, grp combat.ICDGroup, char int) float64 {

	//check if we need to start timer
	if !t.icdDamageTagOnTimer[char][tag] {
		t.icdDamageTagOnTimer[char][tag] = true
		t.ResetDamageCounterAfterDelay(tag, grp, char)
	}

	val := t.icdDamageTagCounter[char][tag]
	//increment the counter
	//if counter > length, then reset back to 0
	t.icdDamageTagCounter[char][tag]++
	if t.icdDamageTagCounter[char][tag] == len(combat.ICDGroupDamageSequence[grp]) {
		t.icdDamageTagCounter[char][tag] = 0
	}

	//true if group seq is 1
	return combat.ICDGroupDamageSequence[grp][val]
}

func (t *LeaLotus) ResetDamageCounterAfterDelay(tag combat.ICDTag, grp combat.ICDGroup, char int) {
	t.Core.Tasks.Add(func() {
		//set the counter back to 0
		t.icdDamageTagCounter[char][tag] = 0
		t.icdDamageTagOnTimer[char][tag] = false
		t.Core.Log.NewEvent("damage counter reset", glog.LogICDEvent, char).
			Write("tag", tag).
			Write("grp", grp)
	}, combat.ICDGroupResetTimer[grp]-1)
	t.Core.Log.NewEvent("damage reset timer set", glog.LogICDEvent, char).
		Write("tag", tag).
		Write("grp", grp).
		Write("reset", t.Core.F+combat.ICDGroupResetTimer[grp]-1)
}

func (t *LeaLotus) ResetTagCounterAfterDelay(tag combat.ICDTag, grp combat.ICDGroup, char int) {
	t.Core.Tasks.Add(func() {
		//set the counter back to 0
		t.icdTagCounter[char][tag] = 0
		t.icdTagOnTimer[char][tag] = false
		t.Core.Log.NewEvent("ele app counter reset", glog.LogICDEvent, char).
			Write("tag", tag).
			Write("grp", grp)
	}, combat.ICDGroupResetTimer[grp]-1)
	t.Core.Log.NewEvent("ele app reset timer set", glog.LogICDEvent, char).
		Write("tag", tag).
		Write("grp", grp).
		Write("reset", t.Core.F+combat.ICDGroupResetTimer[grp]-1)
}
