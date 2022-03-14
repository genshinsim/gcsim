package fischl

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func (c *char) a4() {
	last := 0
	cb := func(args ...interface{}) bool {

		t := args[0].(coretype.Target)
		ae := args[1].(*coretype.AttackEvent)

		if ae.Info.ActorIndex != c.Core.ActiveChar {
			return false
		}
		//do nothing if oz not on field
		if c.ozActiveUntil < c.Core.Frame {
			return false
		}
		if c.Core.Frame-30 < last && last != 0 {
			return false
		}
		last = c.Core.Frame

		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Fischl A4",
			AttackTag:  core.AttackTagElementalArt,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupFischl,
			StrikeType: core.StrikeTypePierce,
			Element:    core.Electro,
			Durability: 25,
			Mult:       0.8,
		}
		// TODO: Ugly hack needed to maintain snapshot logs...
		// Technically should have a separate snapshot for each attack info?
		// ai.ModsLog = c.ozSnapshot.Info.ModsLog
		// A4 uses Oz Snapshot
		c.Core.Combat.QueueAttackWithSnap(ai, c.ozSnapshot.Snapshot, core.NewDefSingleTarget(t.Index(), coretype.TargettableEnemy), 0)

		return false
	}
	c.Core.Subscribe(core.OnOverload, cb, "fischl-a4")
	c.Core.Subscribe(core.OnElectroCharged, cb, "fischl-a4")
	c.Core.Subscribe(core.OnSuperconduct, cb, "fischl-a4")
	c.Core.Subscribe(coretype.OnSwirlElectro, cb, "fischl-a4")
	c.Core.Subscribe(core.OnCrystallizeElectro, cb, "fischl-a4")
}
