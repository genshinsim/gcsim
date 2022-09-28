package reactable

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

func TestHydroBloom(t *testing.T) {
	c := testCore()
	trg := addTargetToCore(c)
	c.Init()

	trg.AttachOrRefill(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 25,
		},
	})
	next := &combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 50,
		},
	}
	trg.React(next)
	// should create a seed, explodes after 5s
	advanceCoreFrame(c)

	if c.Combat.GadgetCount() != 1 {
		t.Errorf("expected created a gadget (bloom), got %v", c.Combat.GadgetCount())
	}
	if !next.Reacted {
		t.Errorf("expected reacted to be true, got %v", next.Reacted)
	}
	if trg.AuraContains(attributes.Hydro, attributes.Dendro) {
		t.Error("expecting target to not contain any remaining hydro or dendro aura")
	}
}

func TestDendroBloom(t *testing.T) {
	c := testCore()
	trg := addTargetToCore(c)
	c.Init()

	trg.AttachOrRefill(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 50,
		},
	})
	next := &combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 25,
		},
	}
	trg.React(next)
	// should create a seed, explodes after 5s
	advanceCoreFrame(c)

	if c.Combat.GadgetCount() != 1 {
		t.Errorf("expected created a gadget (bloom), got %v", c.Combat.GadgetCount())
	}
	if !next.Reacted {
		t.Errorf("expected reacted to be true, got %v", next.Reacted)
	}
	if trg.AuraContains(attributes.Hydro, attributes.Dendro) {
		t.Error("expecting target to not contain any remaining hydro or dendro aura")
	}
}

// testing if it could create 2 seeds with dendro -> ec
func TestECBloom(t *testing.T) {
	c := testCore()
	trg := addTargetToCore(c)
	c.Init()

	trg.AttachOrRefill(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 25,
		},
	})
	trg.React(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Electro,
			Durability: 25,
		},
	})
	// reduce a bit ec aura
	for i := 0; i < 10; i++ {
		advanceCoreFrame(c)
	}

	// dendro -> ec
	next := &combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 25,
		},
	}
	trg.React(next)

	// t.Logf("active auras %v", trg.ActiveAuraString())
	if c.Combat.GadgetCount() != 2 {
		t.Errorf("expected created 2 blooms, got %v", c.Combat.GadgetCount())
	}
}
