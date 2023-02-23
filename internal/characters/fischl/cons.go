package fischl

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (c *char) c6() {
	//this is on attack animation state, not attack landed
	//TODO: this used to be on PostAttack, changed to OnAttack
	//i think this might be more accurate to be OnAttackWillLand? or on animation state change?
	c.Core.Events.Subscribe(event.OnAttack, func(_ ...interface{}) bool {
		//do nothing if oz not on field
		if c.ozActiveUntil < c.Core.F {
			return false
		}
		c.Core.Log.NewEvent("fischl c6 triggered", glog.LogCharacterEvent, c.Index)
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Fischl C6",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagElementalArt,
			ICDGroup:   attacks.ICDGroupFischl,
			StrikeType: attacks.StrikeTypePierce,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       0.3,
		}
		// TODO: Ugly hack needed to maintain snapshot logs...
		// Technically should have a separate snapshot for each attack info?
		// ai.ModsLog = c.ozSnapshot.Info.ModsLog
		// C4 uses Oz Snapshot
		c.Core.QueueAttackWithSnap(
			ai,
			c.ozSnapshot.Snapshot,
			combat.NewBoxHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				combat.Point{Y: -1},
				0.1,
				1,
			),
			c.ozTravel,
		)
		return false
	}, "fischl-c6")
}
