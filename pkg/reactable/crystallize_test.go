package reactable

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func TestCrystallizeCryo(t *testing.T) {
	c := testCore()
	trg := addTargetToCore(c)

	ok := false
	c.Events.Subscribe(event.OnCrystallizeCryo, func(args ...any) bool {
		ok = true
		return true
	}, "crystallize-check")

	c.Init()

	trg.AttachOrRefill(&info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Cryo,
			Durability: 25,
		},
	})
	trg.React(&info.AttackEvent{
		Info: info.AttackInfo{
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

	if !durApproxEqual(7.5, trg.Durability[info.ReactionModKeyCryo], 0.0001) {
		t.Errorf("expecting 7.5 pyro left, got %v", trg.Durability[info.ReactionModKeyCryo])
	}
}

func TestCrystallizePyro(t *testing.T) {
	c := testCore()
	trg := addTargetToCore(c)

	ok := false
	c.Events.Subscribe(event.OnCrystallizePyro, func(args ...any) bool {
		ok = true
		return true
	}, "crystallize-check")

	c.Init()

	trg.AttachOrRefill(&info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Pyro,
			Durability: 25,
		},
	})
	// force on burning
	trg.Durability[info.ReactionModKeyBurning] = 50

	trg.React(&info.AttackEvent{
		Info: info.AttackInfo{
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
	if !durApproxEqual(7.5, trg.Durability[info.ReactionModKeyPyro], 0.0001) {
		t.Errorf("expecting 7.5 pyro left, got %v", trg.Durability[info.ReactionModKeyPyro])
	}
	if !durApproxEqual(37.5, trg.Durability[info.ReactionModKeyBurning], 0.0001) {
		t.Errorf("expecting 37.5 burning left, got %v", trg.Durability[info.ReactionModKeyBurning])
	}
}
