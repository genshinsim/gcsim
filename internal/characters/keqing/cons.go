package keqing

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (c *char) c2() {
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if c.Core.F < c.c2ICD {
			return false
		}
		if c.Core.Rand.Float64() < 0.5 {
			c.c2ICD = c.Core.F + 300
			c.Core.QueueParticle("keqing", 1, attributes.Electro, 100)
			c.Core.Log.NewEvent("keqing c2 proc'd", glog.LogCharacterEvent, c.Index).
				Write("next ready", c.c2ICD)
		}
		return false
	}, "keqing-c2")
}

func (c *char) c4() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.25

	cb := func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}

		c.AddStatMod("keqing-c4", 600, attributes.ATKP, func() ([]float64, bool) {
			return m, true
		})

		return false
	}

	c.Core.Events.Subscribe(event.OnOverload, cb, "keqing-c4")
	c.Core.Events.Subscribe(event.OnElectroCharged, cb, "keqing-c4")
	c.Core.Events.Subscribe(event.OnSuperconduct, cb, "keqing-c4")
	c.Core.Events.Subscribe(event.OnSwirlElectro, cb, "keqing-c4")
	c.Core.Events.Subscribe(event.OnCrystallizeElectro, cb, "keqing-c4")
}

func (c *char) c6(src string) {
	m := make([]float64, attributes.EndStatType)
	m[attributes.ElectroP] = 0.06

	c.AddStatMod("keqing-c6-"+src, 480, attributes.ElectroP, func() ([]float64, bool) {
		return m, true
	})
}
