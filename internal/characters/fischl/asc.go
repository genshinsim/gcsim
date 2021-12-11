package fischl

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) a4() {
	last := 0
	cb := func(args ...interface{}) bool {

		t := args[0].(core.Target)
		ae := args[1].(*core.AttackEvent)

		if ae.Info.ActorIndex != c.Core.ActiveChar {
			return false
		}
		//do nothing if oz not on field
		if c.ozActiveUntil < c.Core.F {
			return false
		}
		if c.Core.F-30 < last && last != 0 {
			return false
		}
		last = c.Core.F
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
		c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(t.Index(), core.TargettableEnemy), 0, 1)

		return false
	}
	c.Core.Events.Subscribe(core.OnOverload, cb, "fischl-a4")
	c.Core.Events.Subscribe(core.OnElectroCharged, cb, "fischl-a4")
	c.Core.Events.Subscribe(core.OnSuperconduct, cb, "fischl-a4")
	c.Core.Events.Subscribe(core.OnSwirlElectro, cb, "fischl-a4")
	c.Core.Events.Subscribe(core.OnCrystallizeElectro, cb, "fischl-a4")
}
