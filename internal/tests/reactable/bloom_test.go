package reactable_test

import (
	"math"
	"testing"

	"github.com/genshinsim/gcsim/internal/template/dendrocore"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func TestHydroBloom(t *testing.T) {
	c, trg := makeCore(1)
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}

	c.QueueAttackEvent(&info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 25,
		},
		Pattern: combat.NewCircleHitOnTarget(info.Point{}, nil, 100),
	}, 0)
	advanceCoreFrame(c)

	c.QueueAttackEvent(&info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 50,
		},
		Pattern: combat.NewCircleHitOnTarget(info.Point{}, nil, 100),
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
}

func TestDendroBloom(t *testing.T) {
	c, trg := makeCore(1)
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}

	c.QueueAttackEvent(&info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 50,
		},
		Pattern: combat.NewCircleHitOnTarget(info.Point{}, nil, 100),
	}, 0)
	advanceCoreFrame(c)

	c.QueueAttackEvent(&info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 25,
		},
		Pattern: combat.NewCircleHitOnTarget(info.Point{}, nil, 100),
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
}

func TestECBloom(t *testing.T) {
	c, trg := makeCore(1)
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}

	// create 2 seeds with ec
	c.QueueAttackEvent(&info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 50,
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
	advanceCoreFrame(c)

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
			Element:    attributes.Dendro,
			Durability: 25,
		},
		Pattern: combat.NewCircleHitOnTarget(trg[0], nil, 100),
	}, 0)
	advanceCoreFrame(c)

	for range dendrocore.Delay + 1 {
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
	for range 6 {
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
				Element:    attributes.Dendro,
				Durability: 25,
			},
			Pattern: combat.NewCircleHitOnTarget(trg[0], nil, 100),
		}, 0)
		advanceCoreFrame(c)
	}

	for range dendrocore.Delay + 1 {
		advanceCoreFrame(c)
	}

	if c.Combat.GadgetCount() != 5 {
		t.Errorf("expected only 5 seeds remaining, got %v", c.Combat.GadgetCount())
	}

	// find oldest
	f := math.MaxInt
	oldest := -1
	for i, v := range c.Combat.Gadgets() {
		if v == nil || v.GadgetTyp() != info.GadgetTypDendroCore {
			continue
		}
		if v.Src() < f {
			f = v.Src()
			oldest = i
		}
	}
	og := c.Combat.Gadget(oldest)
	if og.Src() != 3+dendrocore.Delay {
		t.Errorf("expecting oldest gadget to be from frame %v, got %v", dendrocore.Delay+3, og.Src())
	}
}
