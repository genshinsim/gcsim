package dendrocore

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/gadget"
)

const (
	BloomMultiplier      = 2
	BurgeonMultiplier    = 3
	HyperbloomMultiplier = 3
	Delay                = 30
)

type Gadget struct {
	*gadget.Gadget
	srcFrame  int
	CharIndex int
}

func New(c *core.Core, shp info.Shape, a *info.AttackEvent) *Gadget {
	s := &Gadget{
		srcFrame:  c.F,
		CharIndex: a.Info.ActorIndex,
	}

	circ, ok := shp.(*info.Circle)
	if !ok {
		panic("rectangle target hurtbox is not supported for dendro core spawning")
	}

	// for simplicity, seeds spawn randomly at radius + 0.5
	r := circ.Radius() + 0.5
	s.Gadget = gadget.New(c, info.CalcRandomPointFromCenter(circ.Pos(), r, r, c.Rand), 2, info.GadgetTypDendroCore)
	s.Duration = 300 // ??

	char := s.Core.Player.ByIndex(a.Info.ActorIndex)

	explode := func(reason string) func() {
		return func() {
			s.Core.Tasks.Add(func() {
				ai, snap := NewBloomAttack(char, s, nil)
				ap := combat.NewCircleHitOnTarget(s, nil, 5)
				c.QueueAttackWithSnap(ai, snap, ap, 0)

				// self damage
				ai.Abil += info.SelfDamageSuffix
				ai.FlatDmg = 0.05 * ai.FlatDmg
				ap.SkipTargets[info.TargettablePlayer] = false
				ap.SkipTargets[info.TargettableEnemy] = true
				ap.SkipTargets[info.TargettableGadget] = true
				c.QueueAttackWithSnap(ai, snap, ap, 0)

				c.Log.NewEvent(
					"dendro core "+reason,
					glog.LogElementEvent,
					char.Index(),
				).Write("src", s.Src())
			}, 1)
		}
	}
	s.OnExpiry = explode("expired")
	s.OnKill = explode("killed")

	return s
}

func (s *Gadget) Tick() {
	// this is needed since gadget tick
	s.Gadget.Tick()
}

func (s *Gadget) HandleAttack(atk *info.AttackEvent) float64 {
	s.Core.Events.Emit(event.OnGadgetHit, s, atk)
	s.Attack(atk, nil)
	return 0
}

func (s *Gadget) Attack(atk *info.AttackEvent, evt glog.Event) (float64, bool) {
	if atk.Info.Durability < info.ZeroDur {
		return 0, false
	}

	char := s.Core.Player.ByIndex(atk.Info.ActorIndex)
	// only contact with pyro/electro to trigger burgeon/hyperbloom accordingly
	switch atk.Info.Element {
	case attributes.Electro:
		// trigger hyperbloom targets the nearest enemy
		// it can also do damage to player in small aoe
		s.Core.Tasks.Add(func() {
			ai, snap := NewHyperbloomAttack(char, s)
			// queue dmg nearest enemy within radius 15
			enemy := s.Core.Combat.ClosestEnemyWithinArea(combat.NewCircleHitOnTarget(s.Gadget, nil, 15), nil)
			if enemy != nil {
				ap := combat.NewCircleHitOnTarget(enemy, nil, 1)
				s.Core.QueueAttackWithSnap(ai, snap, ap, 0)

				// also queue self damage
				ai.Abil += info.SelfDamageSuffix
				ai.FlatDmg = 0.05 * ai.FlatDmg
				ap.SkipTargets[info.TargettablePlayer] = false
				ap.SkipTargets[info.TargettableEnemy] = true
				ap.SkipTargets[info.TargettableGadget] = true
				s.Core.QueueAttackWithSnap(ai, snap, ap, 0)
			}
		}, 60)

		s.OnKill = nil
		s.Kill()
		s.Core.Events.Emit(event.OnHyperbloom, s, atk)
		s.Core.Log.NewEvent(
			"hyperbloom triggered",
			glog.LogElementEvent,
			char.Index(),
		).
			Write("dendro_core_char", s.CharIndex).
			Write("dendro_core_src", s.Src())
	case attributes.Pyro:
		// trigger burgeon, aoe dendro damage
		// self damage
		s.Core.Tasks.Add(func() {
			ai, snap := NewBurgeonAttack(char, s)
			ap := combat.NewCircleHitOnTarget(s, nil, 5)
			s.Core.QueueAttackWithSnap(ai, snap, ap, 0)

			// queue self damage
			ai.Abil += info.SelfDamageSuffix
			ai.FlatDmg = 0.05 * ai.FlatDmg
			ap.SkipTargets[info.TargettablePlayer] = false
			ap.SkipTargets[info.TargettableEnemy] = true
			ap.SkipTargets[info.TargettableGadget] = true
			s.Core.QueueAttackWithSnap(ai, snap, ap, 0)
		}, 1)

		s.OnKill = nil
		s.Kill()
		s.Core.Events.Emit(event.OnBurgeon, s, atk)
		s.Core.Log.NewEvent(
			"burgeon triggered",
			glog.LogElementEvent,
			char.Index(),
		).
			Write("dendro_core_char", s.CharIndex).
			Write("dendro_core_src", s.Src())
	default:
		return 0, false
	}

	return 0, false
}

func NewBloomAttack(char *character.CharWrapper, src info.Target, modify func(*info.AttackInfo)) (info.AttackInfo, info.Snapshot) {
	em := char.Stat(attributes.EM)
	ai := info.AttackInfo{
		ActorIndex:       char.Index(),
		DamageSrc:        src.Key(),
		Element:          attributes.Dendro,
		AttackTag:        attacks.AttackTagBloom,
		ICDTag:           attacks.ICDTagBloomDamage,
		ICDGroup:         attacks.ICDGroupReactionA,
		StrikeType:       attacks.StrikeTypeDefault,
		Abil:             string(info.ReactionTypeBloom),
		IgnoreDefPercent: 1,
	}
	if modify != nil {
		modify(&ai)
	}
	flatdmg, snap := combat.CalcReactionDmg(char.Base.Level, char, ai, em)
	ai.FlatDmg = BloomMultiplier * flatdmg
	return ai, snap
}

func NewBurgeonAttack(char *character.CharWrapper, src info.Target) (info.AttackInfo, info.Snapshot) {
	em := char.Stat(attributes.EM)
	ai := info.AttackInfo{
		ActorIndex:       char.Index(),
		DamageSrc:        src.Key(),
		Element:          attributes.Dendro,
		AttackTag:        attacks.AttackTagBurgeon,
		ICDTag:           attacks.ICDTagBurgeonDamage,
		ICDGroup:         attacks.ICDGroupReactionA,
		StrikeType:       attacks.StrikeTypeDefault,
		Abil:             string(info.ReactionTypeBurgeon),
		IgnoreDefPercent: 1,
	}
	flatdmg, snap := combat.CalcReactionDmg(char.Base.Level, char, ai, em)
	ai.FlatDmg = BurgeonMultiplier * flatdmg
	return ai, snap
}

func NewHyperbloomAttack(char *character.CharWrapper, src info.Target) (info.AttackInfo, info.Snapshot) {
	em := char.Stat(attributes.EM)
	ai := info.AttackInfo{
		ActorIndex:       char.Index(),
		DamageSrc:        src.Key(),
		Element:          attributes.Dendro,
		AttackTag:        attacks.AttackTagHyperbloom,
		ICDTag:           attacks.ICDTagHyperbloomDamage,
		ICDGroup:         attacks.ICDGroupReactionA,
		StrikeType:       attacks.StrikeTypeDefault,
		Abil:             string(info.ReactionTypeHyperbloom),
		IgnoreDefPercent: 1,
	}
	flatdmg, snap := combat.CalcReactionDmg(char.Base.Level, char, ai, em)
	ai.FlatDmg = HyperbloomMultiplier * flatdmg
	return ai, snap
}

func (s *Gadget) SetDirection(trg info.Point) {}
func (s *Gadget) SetDirectionToClosestEnemy() {}
func (s *Gadget) CalcTempDirection(trg info.Point) info.Point {
	return info.DefaultDirection()
}
