package fischl

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func (c *char) c6() {
	//this is on attack animation state, not attack landed
	c.Core.Subscribe(core.PostAttack, func(args ...interface{}) bool {
		//do nothing if oz not on field
		if c.ozActiveUntil < c.Core.Frame {
			return false
		}
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Fischl C6",
			AttackTag:  core.AttackTagElementalArt,
			ICDTag:     core.ICDTagElementalArt,
			ICDGroup:   core.ICDGroupFischl,
			StrikeType: core.StrikeTypePierce,
			Element:    core.Electro,
			Durability: 25,
			Mult:       0.3,
		}
		// TODO: Ugly hack needed to maintain snapshot logs...
		// Technically should have a separate snapshot for each attack info?
		// ai.ModsLog = c.ozSnapshot.Info.ModsLog
		// C4 uses Oz Snapshot
		c.Core.Combat.QueueAttackWithSnap(ai, c.ozSnapshot.Snapshot, core.NewDefSingleTarget(1, coretype.TargettableEnemy), 0)
		return false
	}, "fischl-c6")
}
