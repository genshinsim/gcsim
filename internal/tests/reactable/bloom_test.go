package reactable_test

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

func TestHydroBloom(t *testing.T) {
	c, trg := makeCore(1)
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}

	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 25,
		},
		Pattern: combat.NewCircleHit(combat.NewCircle(0, 0, 1), 100, false, combat.TargettableEnemy),
	}, 0)
	advanceCoreFrame(c)

	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 50,
		},
		Pattern: combat.NewCircleHit(combat.NewCircle(0, 0, 1), 100, false, combat.TargettableEnemy),
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
}

func TestDendroBloom(t *testing.T) {
	c, trg := makeCore(1)
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}

	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 50,
		},
		Pattern: combat.NewCircleHit(combat.NewCircle(0, 0, 1), 100, false, combat.TargettableEnemy),
	}, 0)
	advanceCoreFrame(c)

	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 25,
		},
		Pattern: combat.NewCircleHit(combat.NewCircle(0, 0, 1), 100, false, combat.TargettableEnemy),
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
}

func TestECBloom(t *testing.T) {
	c, trg := makeCore(1)
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}

	// create 2 seeds with ec
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 50,
		},
		Pattern: combat.NewCircleHit(trg[0], 100, false, combat.TargettableEnemy, combat.TargettableGadget),
	}, 0)
	advanceCoreFrame(c)
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Electro,
			Durability: 25,
		},
		Pattern: combat.NewCircleHit(trg[0], 100, false, combat.TargettableEnemy, combat.TargettableGadget),
	}, 0)
	advanceCoreFrame(c)

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
		t.Errorf("expected created 2 gadget (bloom), got %v", c.Combat.GadgetCount())
	}
	if trg[0].AuraContains(attributes.Hydro, attributes.Dendro) {
		t.Errorf("expecting target to not contain any remaining hydro or dendro aura, got %v", trg[0].ActiveAuraString())
	}
}

func TestBloomSeedLimit(t *testing.T) {
	c, trg := makeCore(1)
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}

	for i := 0; i < 6; i++ {
		c.QueueAttackEvent(&combat.AttackEvent{
			Info: combat.AttackInfo{
				Element:    attributes.Hydro,
				Durability: 25,
			},
			Pattern: combat.NewCircleHit(trg[0], 100, false, combat.TargettableEnemy, combat.TargettableGadget),
		}, 0)
		advanceCoreFrame(c)
		c.QueueAttackEvent(&combat.AttackEvent{
			Info: combat.AttackInfo{
				Element:    attributes.Dendro,
				Durability: 25,
			},
			Pattern: combat.NewCircleHit(trg[0], 100, false, combat.TargettableEnemy, combat.TargettableGadget),
		}, 0)
		advanceCoreFrame(c)
	}

	for i := 0; i < reactable.DendroCoreDelay+1; i++ {
		advanceCoreFrame(c)
	}

	if c.Combat.GadgetCount() != 5 {
		t.Errorf("expected only 5 seeds remaining, got %v", c.Combat.GadgetCount())
	}

	dur := getOldestSeedDuration(c)
	if dur != 290 { // oldest dendro core should have duration 290 here
		t.Errorf("expected duration to be %v (EG: first dendro core destroyed), got %v", 290, dur)
	}
}

func getOldestSeedDuration(c *core.Core) int {
	// need to count gadgets for which bloom got destroyed
	for i := 0; i < c.Combat.GadgetCount(); i++ {
		d, ok := c.Combat.Gadget(i).(*reactable.DendroCore)
		if !ok {
			continue
		}
		return d.Duration
	}
	return -1
}
