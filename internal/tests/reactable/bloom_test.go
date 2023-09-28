package reactable_test

import (
	"math"
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
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
		Pattern: combat.NewCircleHitOnTarget(geometry.Point{}, nil, 100),
	}, 0)
	advanceCoreFrame(c)

	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 50,
		},
		Pattern: combat.NewCircleHitOnTarget(geometry.Point{}, nil, 100),
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
		Pattern: combat.NewCircleHitOnTarget(geometry.Point{}, nil, 100),
	}, 0)
	advanceCoreFrame(c)

	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 25,
		},
		Pattern: combat.NewCircleHitOnTarget(geometry.Point{}, nil, 100),
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
	advanceCoreFrame(c)

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
		t.Errorf("expected created 2 gadget (bloom), got %v", c.Combat.GadgetCount())
	}
	if trg[0].AuraContains(attributes.Hydro, attributes.Dendro) {
		t.Errorf("expecting target to not contain any remaining hydro or dendro aura, got %v", trg[0].ActiveAuraString())
	}
}

func TestBloomSeedLimit(t *testing.T) {
	c, trg := makeCore(10)
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}

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
			Element:    attributes.Dendro,
			Durability: 25,
		},
		Pattern: combat.NewCircleHitOnTarget(trg[0], nil, 100),
	}, 0)
	advanceCoreFrame(c)

	for i := 0; i < reactable.DendroCoreDelay+1; i++ {
		advanceCoreFrame(c)
	}

	if c.Combat.GadgetCount() != 5 {
		t.Errorf("expected only 5 seeds remaining, got %v", c.Combat.GadgetCount())
	}
}

func TestBloomOldestDeleted(t *testing.T) {
	c, trg := makeCore(1)
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}

	// oldest should be the 2nd one, which is frame 3 ?
	for i := 0; i < 6; i++ {
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
				Element:    attributes.Dendro,
				Durability: 25,
			},
			Pattern: combat.NewCircleHitOnTarget(trg[0], nil, 100),
		}, 0)
		advanceCoreFrame(c)
	}

	for i := 0; i < reactable.DendroCoreDelay+1; i++ {
		advanceCoreFrame(c)
	}

	if c.Combat.GadgetCount() != 5 {
		t.Errorf("expected only 5 seeds remaining, got %v", c.Combat.GadgetCount())
	}

	// find oldest
	f := math.MaxInt
	oldest := -1
	for i, v := range c.Combat.Gadgets() {
		if v == nil || v.GadgetTyp() != combat.GadgetTypDendroCore {
			continue
		}
		if v.Src() < f {
			f = v.Src()
			oldest = i
		}
	}
	og := c.Combat.Gadget(oldest)
	if og.Src() != 3+reactable.DendroCoreDelay {
		t.Errorf("expecting oldest gadget to be from frame %v, got %v", reactable.DendroCoreDelay+3, og.Src())
	}
}
