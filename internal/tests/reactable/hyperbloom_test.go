package reactable_test

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

func TestHyperbloom(t *testing.T) {
	c, trg := makeCore(2)
	trg[0].SetPos(1, 0)
	trg[1].SetPos(3, 0)
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}
	count := 0
	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		trg := args[0].(combat.Target)
		ae := args[1].(*combat.AttackEvent)
		if trg.Type() == combat.TargettableEnemy && ae.Info.Abil == "hyperbloom" {
			count++
		}
		return false
	}, "hyperbloom")

	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 25,
		},
		Pattern: combat.NewDefSingleTarget(trg[0].Key(), combat.TargettableEnemy),
	}, 0)
	advanceCoreFrame(c)

	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 50,
		},
		Pattern: combat.NewDefSingleTarget(trg[0].Key(), combat.TargettableEnemy),
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
			Element:    attributes.Electro,
			Durability: 25,
		},
		Pattern: combat.NewCircleHit(trg[0], 10, false, combat.TargettableEnemy, combat.TargettableGadget),
	}, 0)

	for i := 0; i < 10; i++ {
		advanceCoreFrame(c)
	}

	if count != 1 { //target 2 should be too far away
		t.Errorf("expecting 1 instance of hyperbloom dmg, got %v", count)
	}
	if c.Combat.GadgetCount() != 0 {
		t.Errorf("gadget should be removed, got %v", c.Combat.GadgetCount())
	}
}

func TestECHyperbloom(t *testing.T) {
	c, trg := makeCore(2)
	trg[0].SetPos(1, 0)
	trg[1].SetPos(3, 0)
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}

	count := 0
	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		trg := args[0].(combat.Target)
		ae := args[1].(*combat.AttackEvent)
		if trg.Type() == combat.TargettableEnemy && ae.Info.Abil == "hyperbloom" {
			count++
		}
		return false
	}, "hyperbloom")

	//create 2 seeds with ec
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 25,
		},
		Pattern: combat.NewDefSingleTarget(trg[0].Key(), combat.TargettableEnemy),
	}, 0)
	advanceCoreFrame(c)
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Electro,
			Durability: 25,
		},
		Pattern: combat.NewDefSingleTarget(trg[0].Key(), combat.TargettableEnemy),
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
		Pattern: combat.NewCircleHit(trg[0], 100, false, combat.TargettableEnemy, combat.TargettableGadget),
	}, 0)

	for i := 0; i < reactable.DendroCoreDelay+1; i++ {
		advanceCoreFrame(c)
	}

	if c.Combat.GadgetCount() != 2 {
		t.Errorf("expected 2 bloom gadgets, got %v", c.Combat.GadgetCount())
	}

	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Electro,
			Durability: 25,
		},
		Pattern: combat.NewCircleHit(trg[0], 100, false, combat.TargettableEnemy, combat.TargettableGadget),
	}, 0)

	for i := 0; i < 10; i++ {
		advanceCoreFrame(c)
	}

	if c.Combat.GadgetCount() != 0 {
		t.Errorf("expected bloom gadgets to be cleared, got %v", c.Combat.GadgetCount())
	}

	if count != 2 {
		t.Errorf("expected 2 instance of hyperbloom damage, got %v", count)
	}

}
