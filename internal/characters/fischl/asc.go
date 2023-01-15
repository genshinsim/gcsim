package fischl

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/gadget"
)

func (c *char) a4() {
	last := 0

	// Hyperbloom comes from a gadget so it doesn't ignore gadgets
	a4cb := func(args ...interface{}) bool {

		ae := args[1].(*combat.AttackEvent)

		if ae.Info.ActorIndex != c.Core.Player.Active() {
			return false
		}
		// do nothing if oz not on field
		if c.ozActiveUntil < c.Core.F {
			return false
		}
		if c.Core.F-30 < last && last != 0 {
			return false
		}
		last = c.Core.F

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Fischl A4",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupFischl,
			StrikeType: combat.StrikeTypePierce,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       0.8,
		}
		// TODO: Ugly hack needed to maintain snapshot logs...
		// Technically should have a separate snapshot for each attack info?
		// ai.ModsLog = c.ozSnapshot.Info.ModsLog
		// A4 uses Oz Snapshot

		// TODO: this should target closest enemy within 15m of "elemental reaction position"
		c.Core.QueueAttackWithSnap(
			ai,
			c.ozSnapshot.Snapshot,
			combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 0.5),
			3)
		return false
	}

	a4cbNoGadget := func(args ...interface{}) bool {
		if _, ok := args[0].(*gadget.Gadget); ok {
			return false
		}

		return a4cb(args...)
	}

	c.Core.Events.Subscribe(event.OnOverload, a4cbNoGadget, "fischl-a4")
	c.Core.Events.Subscribe(event.OnElectroCharged, a4cbNoGadget, "fischl-a4")
	c.Core.Events.Subscribe(event.OnSuperconduct, a4cbNoGadget, "fischl-a4")
	c.Core.Events.Subscribe(event.OnSwirlElectro, a4cbNoGadget, "fischl-a4")
	c.Core.Events.Subscribe(event.OnCrystallizeElectro, a4cbNoGadget, "fischl-a4")
	c.Core.Events.Subscribe(event.OnHyperbloom, a4cb, "fischl-a4")
	c.Core.Events.Subscribe(event.OnQuicken, a4cbNoGadget, "fischl-a4")
	c.Core.Events.Subscribe(event.OnAggravate, a4cbNoGadget, "fischl-a4")
}
