package reactable

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func TestFreezePlusCryoHydro(t *testing.T) {
	c := testCore()
	trg := addTargetToCore(c)
	c.Init()

	trg.AttachOrRefill(&info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Cryo,
			Durability: 25,
		},
	})
	trg.React(&info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 25,
		},
	})
	// without ticking, we should have 50 frozen here
	if !durApproxEqual(40, trg.Durability[info.ReactionModKeyFrozen], 0.00001) {
		t.Errorf("frozen expected to be 40, got %v", trg.Durability[info.ReactionModKeyFrozen])
		t.FailNow()
	}

	trg.AttachOrRefill(&info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Cryo,
			Durability: 25,
		},
	})

	// should have frozen + cryo here
	if !durApproxEqual(20, trg.Durability[info.ReactionModKeyCryo], 0.00001) {
		t.Errorf("expecting 20 cryo attached, got %v", trg.Durability[info.ReactionModKeyCryo])
	}
}

func TestFreezePlusAddFreeze(t *testing.T) {
	c := testCore()
	trg := addTargetToCore(c)
	c.Init()

	trg.AttachOrRefill(&info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Cryo,
			Durability: 25,
		},
	})
	trg.React(&info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 25,
		},
	})
	// without ticking, we should have 50 frozen here
	if !durApproxEqual(40, trg.Durability[info.ReactionModKeyFrozen], 0.00001) {
		t.Errorf("frozen expected to be 40, got %v", trg.Durability[info.ReactionModKeyFrozen])
		t.FailNow()
	}

	trg.AttachOrRefill(&info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Cryo,
			Durability: 50, // gives us 40 attached
		},
	})
	trg.React(&info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 50,
		},
	})

	// should have frozen + cryo here
	if !durApproxEqual(80, trg.Durability[info.ReactionModKeyFrozen], 0.00001) {
		t.Errorf("expecting 80 frozen attached, got %v", trg.Durability[info.ReactionModKeyFrozen])
	}
}
