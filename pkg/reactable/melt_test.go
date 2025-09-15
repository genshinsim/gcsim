package reactable

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func TestCryoMelt(t *testing.T) {
	c := testCore()
	trg := addTargetToCore(c)
	c.Init()

	trg.AttachOrRefill(&info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Pyro,
			Durability: 25,
		},
	})
	trg.Tick()
	next := &info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Cryo,
			Durability: 50,
		},
	}
	trg.React(next)
	trg.AttachOrRefill(next)
	// check to see if src has amped flag now

	if next.Info.Amped != true && next.Info.AmpMult != 1.5 {
		t.Errorf("expecting amped to be true with factor 1.5, got %v: %v", next.Info.Amped, next.Info.AmpMult)
	}
	if trg.AuraContains(attributes.Hydro, attributes.Pyro) {
		t.Error("expecting target to not contain any remaining cryo or pyro aura")
	}
}

func TestPyroMelt(t *testing.T) {
	c := testCore()
	trg := addTargetToCore(c)
	c.Init()

	trg.AttachOrRefill(&info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Cryo,
			Durability: 25,
		},
	})
	trg.Tick()
	next := &info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Pyro,
			Durability: 50,
		},
	}
	trg.React(next)
	trg.AttachOrRefill(next)
	// check to see if src has amped flag now

	if next.Info.Amped != true && next.Info.AmpMult != 2 {
		t.Errorf("expecting amped to be true with factor 2, got %v: %v", next.Info.Amped, next.Info.AmpMult)
	}
	if trg.AuraContains(attributes.Cryo, attributes.Pyro) {
		t.Error("expecting target to not contain any remaining hydro or pyro aura")
	}
}

func TestPyroFrozenCryoMelt(t *testing.T) {
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
	trg.AttachOrRefill(&info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Cryo,
			Durability: 25,
		},
	})
	trg.Tick()
	next := &info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Pyro,
			Durability: 50,
		},
	}
	trg.React(next)
	trg.AttachOrRefill(next)
	// check to see if src has amped flag now

	if next.Info.Amped != true && next.Info.AmpMult != 2 {
		t.Errorf("expecting amped to be true with factor 2, got %v: %v", next.Info.Amped, next.Info.AmpMult)
	}
	if trg.AuraContains(attributes.Cryo, attributes.Pyro, attributes.Frozen) {
		t.Error("expecting target to not contain any remaining hydro or pyro aura")
	}
}
