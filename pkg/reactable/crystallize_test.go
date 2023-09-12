package reactable

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

func TestCrystallizeCryo(t *testing.T) {

	c := testCore()
	trg := addTargetToCore(c)

	ok := false
	c.Events.Subscribe(event.OnCrystallizeCryo, func(args ...interface{}) bool {
		ok = true
		return true
	}, "crystallize-check")

	c.Init()

	trg.AttachOrRefill(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Cryo,
			Durability: 25,
		},
	})
	trg.React(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Geo,
			Durability: 25,
		},
	})

	// check shield
	if !ok {
		t.Errorf("expecting crystallize to have occured")
		t.FailNow()
	}
	if trg.core.Player.Shields.Count() == 0 {
		t.Errorf("expecting player to be shielded")
	}

	if !durApproxEqual(7.5, trg.Durability[ModifierCryo], 0.0001) {
		t.Errorf("expecting 7.5 pyro left, got %v", trg.Durability[ModifierCryo])
	}
}

func TestCrystallizePyro(t *testing.T) {

	c := testCore()
	trg := addTargetToCore(c)

	ok := false
	c.Events.Subscribe(event.OnCrystallizePyro, func(args ...interface{}) bool {
		ok = true
		return true
	}, "crystallize-check")

	c.Init()

	trg.AttachOrRefill(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Pyro,
			Durability: 25,
		},
	})
	// force on burning
	trg.Durability[ModifierBurning] = 50

	trg.React(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Geo,
			Durability: 25,
		},
	})

	// check shield
	if !ok {
		t.Errorf("expecting crystallize to have occured")
		t.FailNow()
	}
	if trg.core.Player.Shields.Count() == 0 {
		t.Errorf("expecting player to be shielded")
	}
	if !durApproxEqual(7.5, trg.Durability[ModifierPyro], 0.0001) {
		t.Errorf("expecting 7.5 pyro left, got %v", trg.Durability[ModifierPyro])
	}
	if !durApproxEqual(37.5, trg.Durability[ModifierBurning], 0.0001) {
		t.Errorf("expecting 37.5 burning left, got %v", trg.Durability[ModifierBurning])
	}
}
