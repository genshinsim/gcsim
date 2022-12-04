package beidou

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

const c4key = "beidou-c4"

func (c *char) c4() {
	c.Core.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
		di := args[0].(player.DrainInfo)
		if !di.External {
			return false
		}
		if di.Amount <= 0 {
			return false
		}
		if c.Core.Player.Active() != c.Index {
			return false
		}
		c.Core.Status.Add("beidouc4", 600)
		c.Core.Log.NewEvent("c4 triggered on damage", glog.LogCharacterEvent, c.Index).
			Write("expiry", c.Core.F+600)
		return false
	}, "beidouc4")

	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
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
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       0.2,
		}
		c.Core.QueueAttack(ai, combat.NewSingleTargetHit(t.Key()), 0, 1)
		return false
	}, "beidou-c4")
}
