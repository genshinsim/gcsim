package beidou

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

const c4key = "beidou-c4"

func (c *char) c4() {
	c.Core.Events.Subscribe(event.OnCharacterHurt, func(_ ...interface{}) bool {
		if c.Core.Player.Active() != c.Index {
			return false
		}
		c.Core.Status.Add("beidouc4", 600)
		c.Core.Log.NewEvent("c4 triggered on damage", glog.LogCharacterEvent, c.Index).
			Write("expiry", c.Core.F+600)
		return false
	}, "beidouc4")

	c.Core.Events.Subscribe(event.OnDamage, func(evt event.EventPayload) bool {
		t := args[0].(combat.Target)
		ae := args[1].(*combat.AttackEvent)
		if ae.Info.ActorIndex != c.Index {
			return false
		}
		if ae.Info.AttackTag != combat.AttackTagNormal && ae.Info.AttackTag != combat.AttackTagExtra {
			return false
		}
		if !c.StatusIsActive(c4key) {
			return false
		}

		c.Core.Log.NewEvent("c4 proc'd on attack", glog.LogCharacterEvent, c.Index).
			Write("char", c.Index)
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Beidou C4",
			AttackTag:  combat.AttackTagNone,
			ICDTag:     combat.ICDTagElementalBurst,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeBlunt,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       0.2,
		}
		c.Core.QueueAttack(ai, combat.NewDefSingleTarget(t.Index(), t.Type()), 0, 1)
		return false
	}, "beidou-c4")
}
