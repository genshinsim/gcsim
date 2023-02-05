package keqing

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/gadget"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c2ICDKey = "keqing-c2-icd"

func (c *char) c2() {
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		e, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagNormal && atk.Info.AttackTag != combat.AttackTagExtra {
			return false
		}
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if c.Core.Player.Active() != c.Index {
			return false
		}
		if !e.AuraContains(attributes.Electro) {
			return false
		}
		if c.StatusIsActive(c2ICDKey) {
			return false
		}
		if c.Core.Rand.Float64() < 0.5 {
			c.AddStatus(c2ICDKey, 300, true)
			c.Core.QueueParticle("keqing-c2", 1, attributes.Electro, c.ParticleDelay)
			c.Core.Log.NewEvent("keqing c2 proc'd", glog.LogCharacterEvent, c.Index)
		}
		return false
	}, "keqing-c2")
}

func (c *char) c4() {
	cb := func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("keqing-c4", 600),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return c.c4buff, true
			},
		})

		return false
	}

	cbNoGadget := func(args ...interface{}) bool {
		if _, ok := args[0].(*gadget.Gadget); ok {
			return false
		}
		return cb(args...)
	}
	c.Core.Events.Subscribe(event.OnOverload, cbNoGadget, "keqing-c4")
	c.Core.Events.Subscribe(event.OnElectroCharged, cbNoGadget, "keqing-c4")
	c.Core.Events.Subscribe(event.OnSuperconduct, cbNoGadget, "keqing-c4")
	c.Core.Events.Subscribe(event.OnSwirlElectro, cbNoGadget, "keqing-c4")
	c.Core.Events.Subscribe(event.OnCrystallizeElectro, cbNoGadget, "keqing-c4")
	c.Core.Events.Subscribe(event.OnQuicken, cbNoGadget, "keqing-c4")
	c.Core.Events.Subscribe(event.OnAggravate, cbNoGadget, "keqing-c4")
	c.Core.Events.Subscribe(event.OnHyperbloom, cb, "keqing-c4")
}

func (c *char) c6(src string) {
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("keqing-c6-"+src, 480),
		AffectedStat: attributes.ElectroP,
		Amount: func() ([]float64, bool) {
			return c.c6buff, true
		},
	})
}
