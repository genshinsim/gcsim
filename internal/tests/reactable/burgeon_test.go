package reactable_test

import (
	"testing"

	"github.com/genshinsim/gcsim/internal/template/dendrocore"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func TestBurgeon(t *testing.T) {
	c, trg := makeCore(2)
	trg[0].SetPos(info.Point{X: 1, Y: 0})
	trg[1].SetPos(info.Point{X: 3, Y: 0})
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}
	count := 0
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...any) bool {
		trg := args[0].(info.Target)
		ae := args[1].(*info.AttackEvent)
		if trg.Type() == info.TargettableEnemy && ae.Info.Abil == "burgeon" {
			count++
		}
		return false
	}, "burgeon")

	c.QueueAttackEvent(&info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 25,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 0)
	advanceCoreFrame(c)

	c.QueueAttackEvent(&info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 50,
		},
		Pattern: combat.NewSingleTargetHit(trg[0].Key()),
	}, 0)

	// should create a seed, explodes after 5s
	for range dendrocore.Delay + 1 {
		advanceCoreFrame(c)
	}
	if c.Combat.GadgetCount() != 1 {
		t.Errorf("expected created a gadget (bloom), got %v", c.Combat.GadgetCount())
	}
	if trg[0].AuraContains(attributes.Hydro, attributes.Dendro) {
		t.Errorf("expecting target to not contain any remaining hydro or dendro aura, got %v", trg[0].ActiveAuraString())
	}

	c.QueueAttackEvent(&info.AttackEvent{
		Info: info.AttackInfo{
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
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...any) bool {
		trg := args[0].(info.Target)
		ae := args[1].(*info.AttackEvent)
		if trg.Type() == info.TargettableEnemy && ae.Info.Abil == "burgeon" {
			count++
		}
		return false
	}, "burgeon")

	// create 2 seeds with ec
	c.QueueAttackEvent(&info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 25,
		},
		Pattern: combat.NewCircleHitOnTarget(trg[0], nil, 100),
	}, 0)
	advanceCoreFrame(c)
	c.QueueAttackEvent(&info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Electro,
			Durability: 25,
		},
		Pattern: combat.NewCircleHitOnTarget(trg[0], nil, 100),
	}, 0)
	// reduce aura a bit
	for range 10 {
		advanceCoreFrame(c)
	}

	c.QueueAttackEvent(&info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 25,
		},
		Pattern: combat.NewCircleHitOnTarget(trg[0], nil, 100),
	}, 0)

	for range dendrocore.Delay + 1 {
		advanceCoreFrame(c)
	}

	if c.Combat.GadgetCount() != 2 {
		t.Errorf("expected 2 bloom gadgets, got %v", c.Combat.GadgetCount())
	}

	c.QueueAttackEvent(&info.AttackEvent{
		Info: info.AttackInfo{
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
