package fischl

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

func (c *char) a4() {
	last := 0
	cb := func(args ...interface{}) bool {

		t := args[0].(combat.Target)
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
		c.Core.QueueAttackWithSnap(ai, c.ozSnapshot.Snapshot, combat.NewDefSingleTarget(t.Index(), combat.TargettableEnemy), 3)

		return false
	}
	c.Core.Events.Subscribe(event.OnOverload, cb, "fischl-a4")
	c.Core.Events.Subscribe(event.OnElectroCharged, cb, "fischl-a4")
	c.Core.Events.Subscribe(event.OnSuperconduct, cb, "fischl-a4")
	c.Core.Events.Subscribe(event.OnSwirlElectro, cb, "fischl-a4")
	c.Core.Events.Subscribe(event.OnCrystallizeElectro, cb, "fischl-a4")
	c.Core.Events.Subscribe(event.OnQuicken, cb, "fischl-a4")
	c.Core.Events.Subscribe(event.OnAggravate, cb, "fischl-a4")
}
