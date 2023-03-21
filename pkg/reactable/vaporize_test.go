package reactable

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

func TestHydroVaporize(t *testing.T) {
	c := testCore()
	trg := addTargetToCore(c)
	c.Init()

	var next *combat.AttackEvent
	c.Events.Subscribe(event.OnVaporize, func(args ...interface{}) bool {
		if ae, ok := args[1].(*combat.AttackEvent); ok {
			next = ae
		}
		return false
	}, "vape-test")

	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Pyro,
			Durability: 50,
		},
		Pattern: combat.NewCircleHitOnTarget(trg, nil, 100),
	}, 0)
	c.Tick()
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 25,
		},
		Pattern: combat.NewCircleHitOnTarget(trg, nil, 100),
	}, 0)
	c.Tick()

	if next == nil {
		t.Error("attack shouldn't be nil!!")
		t.FailNow()
	}

	if next.Info.Amped != true && next.Info.AmpMult != 1.5 {
		t.Errorf("expecting amped to be true with factor 1.5, got %v, %v", next.Info.Amped, next.Info.AmpMult)
	}
	if trg.AuraContains(attributes.Hydro, attributes.Pyro) {
		t.Error("expecting target to not contain any remaining hydro or pyro aura")
	}
}

func TestPyroVaporize(t *testing.T) {
	c := testCore()
	trg := addTargetToCore(c)
	c.Init()
	var next *combat.AttackEvent
	c.Events.Subscribe(event.OnVaporize, func(args ...interface{}) bool {
		if ae, ok := args[1].(*combat.AttackEvent); ok {
			next = ae
		}
		return false
	}, "vape-test")

	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 25,
		},
		Pattern: combat.NewCircleHitOnTarget(trg, nil, 100),
	}, 0)
	c.Tick()
	c.QueueAttackEvent(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Pyro,
			Durability: 50,
		},
		Pattern: combat.NewCircleHitOnTarget(trg, nil, 100),
	}, 0)
	c.Tick()

	if next == nil {
		t.Error("attack shouldn't be nil!!")
		t.FailNow()
	}
	if next.Info.Amped != true && next.Info.AmpMult != 2 {
		t.Errorf("expecting amped to be true with factor 2, got %v: %v", next.Info.Amped, next.Info.AmpMult)
	}
	if trg.AuraContains(attributes.Hydro, attributes.Pyro) {
		t.Error("expecting target to not contain any remaining hydro or pyro aura")
	}
}
