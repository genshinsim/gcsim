package reactable

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

// test hyperbloom interacts with a seed
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
		t.Error("should get one Hyperbloom")
	}
	if src.Info.FlatDmg == 0 {
		t.Error("expected FlatDmg not to be 0")
	}
	if c.Combat.GadgetCount() != 0 {
		t.Errorf("gadget should be removed, got %v", c.Combat.GadgetCount())
	}
}

// hyperbloom with 2 seeds
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

	// should create 2 blooms
	if c.Combat.GadgetCount() != 2 {
		t.Errorf("expected 2 blooms, got %v", c.Combat.GadgetCount())
	}

	var count []*combat.AttackEvent = []*combat.AttackEvent{}
	trg.onDmgCallBack = func(ae *combat.AttackEvent) (float64, bool) {
		if ae.Info.Abil == "Hyperbloom" {
			count = append(count, ae)
		}
		return 0, false
	}

	// queue an aoe electro to proc 2 hyperblooms
	ae := &combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Electro,
			Durability: 50,
		},
	}
	c.QueueAttack(ae.Info, combat.NewCircleHit(trg.self, 10, true, combat.TargettableGadget), -1, 1)
	for i := 0; i < 11; i++ {
		advanceCoreFrame(c)
	}

	// should've cleared blooms
	if c.Combat.GadgetCount() != 0 {
		t.Errorf("expected blooms wiped, still got %v", c.Combat.GadgetCount())
	}
	if len(count) != 2 {
		t.Errorf("Expected 2 hyperblooms, got %v", count)
	}
	for _, v := range count {
		if v.Info.FlatDmg == 0 {
			t.Error("expected FlatDmg not to be 0")
		}
	}
}
