package reactable_test

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

func TestBurgeon(t *testing.T) {
	c, trg := makeCore(2)
	trg[0].SetPos(geometry.Point{X: 1, Y: 0})
	trg[1].SetPos(geometry.Point{X: 3, Y: 0})
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}
	count := 0
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		trg := args[0].(combat.Target)
		ae := args[1].(*combat.AttackEvent)
		if trg.Type() == targets.TargettableEnemy && ae.Info.Abil == "burgeon" {
			count++
		}
		return false
	}, "burgeon")

	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 25,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 0)
	advanceCoreFrame(c)

	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 50,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 0)

	// should create a seed, explodes after 5s
	for i := 0; i < reactable.DendroCoreDelay+1; i++ {
		advanceCoreFrame(c)
	}
	if c.Combat.GadgetCount() != 1 {
		t.Errorf("expected created a gadget (bloom), got %v", c.Combat.GadgetCount())
	}
	if trg[0].AuraContains(attributes.Hydro, attributes.Dendro) {
		t.Errorf("expecting target to not contain any remaining hydro or dendro aura, got %v", trg[0].ActiveAuraString())
	}

	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Pyro,
			Durability: 50,
		},
		Pattern: combat.NewCircleHitOnTarget(trg[0], nil, 10),
	}, 0)
	advanceCoreFrame(c)
	if count != 2 {
		t.Errorf("expecting 2 instance of burgeon dmg, got %v", count)
	}
	if c.Combat.GadgetCount() != 0 {
		t.Errorf("gadget should be removed, got %v", c.Combat.GadgetCount())
	}
}

func TestECBurgeon(t *testing.T) {
	c, trg := makeCore(1)
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}

	count := 0
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		trg := args[0].(combat.Target)
		ae := args[1].(*combat.AttackEvent)
		if trg.Type() == targets.TargettableEnemy && ae.Info.Abil == "burgeon" {
			count++
		}
		return false
	}, "burgeon")

	//create 2 seeds with ec
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 25,
		},
		Pattern: combat.NewCircleHitOnTarget(trg[0], nil, 100),
	}, 0)
	advanceCoreFrame(c)
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Electro,
			Durability: 25,
		},
		Pattern: combat.NewCircleHitOnTarget(trg[0], nil, 100),
	}, 0)
	//reduce aura a bit
	for i := 0; i < 10; i++ {
		advanceCoreFrame(c)
	}

	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 25,
		},
		Pattern: combat.NewCircleHitOnTarget(trg[0], nil, 100),
	}, 0)

	for i := 0; i < reactable.DendroCoreDelay+1; i++ {
		advanceCoreFrame(c)
	}

	if c.Combat.GadgetCount() != 2 {
		t.Errorf("expected 2 bloom gadgets, got %v", c.Combat.GadgetCount())
	}

	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Pyro,
			Durability: 25,
		},
		Pattern: combat.NewCircleHitOnTarget(trg[0], nil, 100),
	}, 0)

	advanceCoreFrame(c)

	if c.Combat.GadgetCount() != 0 {
		t.Errorf("expected bloom gadgets to be cleared, got %v", c.Combat.GadgetCount())
	}

	if count != 2 {
		t.Errorf("expected 2 instance of burgeon damage, got %v", count)
	}

}
