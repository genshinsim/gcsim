package reactable

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

func TestFreezePlusCryoHydro(t *testing.T) {
	c := testCore()
	trg := addTargetToCore(c)
	c.Init()

	trg.AttachOrRefill(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Cryo,
			Durability: 25,
		},
	})
	trg.React(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 25,
		},
	})
	// without ticking, we should have 50 frozen here
	if !durApproxEqual(40, trg.Durability[Frozen], 0.00001) {
		t.Errorf("frozen expected to be 40, got %v", trg.Durability[Frozen])
		t.FailNow()
	}

	trg.AttachOrRefill(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Cryo,
			Durability: 25,
		},
	})

	// should have frozen + cryo here
	if !durApproxEqual(20, trg.Durability[Cryo], 0.00001) {
		t.Errorf("expecting 20 cryo attached, got %v", trg.Durability[Cryo])
	}
}

func TestFreezePlusAddFreeze(t *testing.T) {
	c := testCore()
	trg := addTargetToCore(c)
	c.Init()

	trg.AttachOrRefill(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Cryo,
			Durability: 25,
		},
	})
	trg.React(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 25,
		},
	})
	// without ticking, we should have 50 frozen here
	if !durApproxEqual(40, trg.Durability[Frozen], 0.00001) {
		t.Errorf("frozen expected to be 40, got %v", trg.Durability[Frozen])
		t.FailNow()
	}

	trg.AttachOrRefill(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Cryo,
			Durability: 50, // gives us 40 attached
		},
	})
	trg.React(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 50,
		},
	})

	// should have frozen + cryo here
	if !durApproxEqual(80, trg.Durability[Frozen], 0.00001) {
		t.Errorf("expecting 80 frozen attached, got %v", trg.Durability[Frozen])
	}
}
