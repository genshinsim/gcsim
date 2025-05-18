package fischl

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

const a4IcdKey = "fischl-a4-icd"

// A1 is not implemented:
// TODO: When Fischl hits Oz with a fully-charged Aimed Shot, Oz brings down Thundering Retribution, dealing AoE Electro DMG equal to 152.7% of the arrow's DMG.

// If your current active character triggers an Electro-related Elemental Reaction when Oz is on the field,
// the opponent shall be stricken with Thundering Retribution that deals Electro DMG equal to 80% of Fischl's ATK.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	// Hyperbloom comes from a gadget so it doesn't ignore gadgets
	//nolint:unparam // ignoring for now, event refactor should get rid of bool return of event sub
	a4cb := func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)

		if ae.Info.ActorIndex != c.Core.Player.Active() {
			return false
		}
		// do nothing if oz not on field
		if !c.StatusIsActive(ozActiveKey) {
			return false
		}
		active := c.Core.Player.ActiveChar()
		if active.StatusIsActive(a4IcdKey) {
			return false
		}
		active.AddStatus(a4IcdKey, 0.5*60, true)

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Fischl A4",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupFischl,
			StrikeType: attacks.StrikeTypePierce,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       0.8,
		}

		// A4 uses Oz Snapshot
		// TODO: this should target closest enemy within 15m of "elemental reaction position"
		c.Core.QueueAttackWithSnap(
			ai,
			c.ozSnapshot.Snapshot,
			combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 0.5),
			4)
		return false
	}

	a4cbNoGadget := func(args ...interface{}) bool {
		if _, ok := args[0].(*enemy.Enemy); !ok {
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
