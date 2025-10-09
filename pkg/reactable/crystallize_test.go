package reactable

import (
	"testing"

	"github.com/genshinsim/gcsim/internal/template/crystallize"
	"github.com/genshinsim/gcsim/pkg/core"
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
	advanceCoreFrameMultiple(c, 53)
	count := pickUpCrystallize(c, attributes.Cryo)
	if count > 0 {
		t.Errorf("crystallize shard pickup too early (before 54)")
	}
	advanceCoreFrame(c)
	count = pickUpCrystallize(c, attributes.Cryo)
	if count == 0 {
		t.Errorf("Did not pick up crystallize shard (at 54)")
	}
	// check shield
	if !ok {
		t.Errorf("expecting crystallize to have occured")
		t.FailNow()
	}
	if trg.core.Player.Shields.Count() == 0 {
		t.Errorf("expecting player to be shielded")
	}

	if !durApproxEqual(7.5-0.03508771929824561*54, trg.GetAuraDurability(info.ReactionModKeyCryo), 0.0001) {
		t.Errorf("expecting 7.5 cryo left, got %v", trg.GetAuraDurability(info.ReactionModKeyCryo))
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

	advanceCoreFrameMultiple(c, 53)
	count := pickUpCrystallize(c, attributes.Pyro)
	if count > 0 {
		t.Errorf("crystallize shard pickup too early (before 54)")
	}
	advanceCoreFrame(c)
	count = pickUpCrystallize(c, attributes.Pyro)
	if count == 0 {
		t.Errorf("Did not pick up crystallize shard (at 54)")
	}
	// check shield
	if !ok {
		t.Errorf("expecting crystallize to have occured")
		t.FailNow()
	}
	if trg.core.Player.Shields.Count() == 0 {
		t.Errorf("expecting player to be shielded")
	}
	if !durApproxEqual(7.5-0.03508771929824561*54, trg.GetAuraDurability(info.ReactionModKeyPyro), 0.0001) {
		t.Errorf("expecting 5.605263157894737 pyro left, got %v", trg.GetAuraDurability(info.ReactionModKeyPyro))
	}
	if !durApproxEqual(37.5, trg.GetAuraDurability(info.ReactionModKeyBurning), 0.0001) {
		t.Errorf("expecting 37.5 burning left, got %v", trg.GetAuraDurability(info.ReactionModKeyBurning))
	}
}

func pickUpCrystallize(c *core.Core, pickupEle attributes.Element) int {
	var count int
	for _, g := range c.Combat.Gadgets() {
		shard, ok := g.(*crystallize.Shard)
		// skip if no shard
		if !ok {
			continue
		}
		// skip if shard not specified element
		if pickupEle != attributes.UnknownElement && shard.Shield.Ele != pickupEle {
			continue
		}
		// try to pick up shard and stop if succeeded
		if shard.AddShieldKillShard() {
			count = 1
			break
		}
	}
	return count
}
