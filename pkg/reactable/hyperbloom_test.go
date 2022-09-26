package reactable

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

func TestHyperBloom(t *testing.T) {
	c := testCore()
	trg2 := addTargetToCore(c)
	trg1 := addTargetToCore(c)
	trg2.SetPos(3, 0)
	trg1.SetPos(1, 0)
	c.Init()

	var src *combat.AttackEvent
	trg1.onDmgCallBack = func(atk *combat.AttackEvent) (float64, bool) {
		if atk.Info.Abil == "Hyperbloom" {
			src = atk
		}
		return 0, false
	}

	trg1.AttachOrRefill(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 25,
		},
	})
	trg1.React(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 25,
		},
	})
	// wait a bit for bloom spawns
	for i := 0; i < 10; i++ {
		advanceCoreFrame(c)
	}

	// summon aoe electro to proc hyperbloom nearby trg1
	ae := &combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Electro,
			Durability: 25,
		},
	}
	c.QueueAttack(ae.Info, combat.NewCircleHit(trg2.self, 5, false, combat.TargettableGadget), 1, 1)
	for i := 0; i < 10; i++ {
		advanceCoreFrame(c)
	}

	// trg1 should get hyperbloom damage
	if src == nil {
		t.Error("should get Hyperbloom")
	}
	if c.Combat.GadgetCount() != 0 {
		t.Logf("gadget should be removed, got %v", c.Combat.GadgetCount())
	}
}

func TestECHyperBloom(t *testing.T) {
	c := testCore()
	trg := addTargetToCore(c)
	c.Init()

	// creates 2 seeds with ec
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
	trg.React(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 25,
		},
	})

	var count int = 0
	trg.onDmgCallBack = func(ae *combat.AttackEvent) (float64, bool) {
		if ae.Info.Abil == "Hyperbloom" {
			count++
		}
		return 0, false
	}

	// summon aoe electro to proc 2 hyperbloom
	ae := &combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Electro,
			Durability: 25,
		},
	}
	c.QueueAttack(ae.Info, combat.NewCircleHit(trg.self, 5, false, combat.TargettableGadget), 1, 1)
	for i := 0; i < 10; i++ {
		advanceCoreFrame(c)
	}

	if count != 2 {
		t.Errorf("Expected 2 hyperblooms, got %v", count)
	}
}
