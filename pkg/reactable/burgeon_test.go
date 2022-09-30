package reactable

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

// test burgeon
func TestBurgeon(t *testing.T) {
	c, trg := testCoreWithTrgs(2)
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
		if trg.Type() == combat.TargettableEnemy && ae.Info.Abil == "burgeon" {
			count++
		}
		return false
	}, "burgeon")
	c.QueueAttackEvent(makeSTAttack(attributes.Hydro, 25, 0), 0)
	c.Tick()
	c.QueueAttackEvent(makeSTAttack(attributes.Dendro, 25, 0), 0)

	// wait a bit for bloom spawns
	for i := 0; i < 46; i++ {
		advanceCoreFrame(c)
	}

	if c.Combat.GadgetCount() != 1 {
		t.Errorf("expecting 1 gadget, got %v", c.Combat.GadgetCount())
	}

	// queue aoe pyro to proc burgeon
	ae := &combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Pyro,
			Durability: 25,
		},
	}
	c.QueueAttack(ae.Info, combat.NewCircleHit(trg[1], 5, false, combat.TargettableGadget), 1, 1)
	for i := 0; i < 10; i++ {
		advanceCoreFrame(c)
	}

	if count != 2 {
		t.Errorf("expecting 2 instance of burgeon dmg, got %v", count)
	}
	if c.Combat.GadgetCount() != 0 {
		t.Errorf("gadget should be removed, got %v", c.Combat.GadgetCount())
	}
}

// Burgeon with 2 seeds
// func TestECBurgeon(t *testing.T) {
// 	c := testCore()
// 	trg := addTargetToCore(c)
// 	c.Init()

// 	// creates 2 seeds with ec
// 	trg.AttachOrRefill(&combat.AttackEvent{
// 		Info: combat.AttackInfo{
// 			Element:    attributes.Hydro,
// 			Durability: 25,
// 		},
// 	})
// 	trg.React(&combat.AttackEvent{
// 		Info: combat.AttackInfo{
// 			Element:    attributes.Electro,
// 			Durability: 25,
// 		},
// 	})
// 	// reduce a bit ec aura
// 	for i := 0; i < 10; i++ {
// 		advanceCoreFrame(c)
// 	}
// 	// dendro -> ec
// 	trg.React(&combat.AttackEvent{
// 		Info: combat.AttackInfo{
// 			Element:    attributes.Dendro,
// 			Durability: 25,
// 		},
// 	})
// 	for i := 0; i < 46; i++ {
// 		advanceCoreFrame(c)
// 	}

// 	// should create 2 blooms
// 	if c.Combat.GadgetCount() != 2 {
// 		t.Errorf("expected 2 blooms, got %v", c.Combat.GadgetCount())
// 	}

// 	var count []*combat.AttackEvent = []*combat.AttackEvent{}

// 	// queue an aoe pyro to proc 2 Burgeons
// 	ae := &combat.AttackEvent{
// 		Info: combat.AttackInfo{
// 			Element:    attributes.Pyro,
// 			Durability: 25,
// 		},
// 	}
// 	c.QueueAttack(ae.Info, combat.NewCircleHit(trg.self, 10, true, combat.TargettableGadget), -1, 1)
// 	for i := 0; i < 11; i++ {
// 		advanceCoreFrame(c)
// 	}

// 	// should've cleared blooms
// 	if c.Combat.GadgetCount() != 0 {
// 		t.Errorf("expected blooms wiped, still got %v", c.Combat.GadgetCount())
// 	}
// 	if len(count) != 2 {
// 		t.Errorf("Expected 2 Burgeons, got %v", count)
// 	}
// 	for _, v := range count {
// 		if v.Info.FlatDmg == 0 {
// 			t.Error("expected FlatDmg not to be 0")
// 		}
// 	}
// }
