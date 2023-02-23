package beidou

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

const c4Key = "beidou-c4"

// Upon being attacked, Beidou's Normal Attacks gain an additional instance of 20% Electro DMG for 10s.
func (c *char) c4Init() {
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
		c.Core.Status.Add(c4Key, 600)
		c.Core.Log.NewEvent("c4 triggered on damage", glog.LogCharacterEvent, c.Index).
			Write("expiry", c.Core.F+600)
		return false
	}, c4Key)
}

// TODO: this should also be added to her CA
// Beidou's Normal Attacks gain an additional instance of 20% Electro DMG for 10s.
func (c *char) makeC4Callback() combat.AttackCBFunc {
	if c.Base.Cons < 4 {
		return nil
	}
	return func(a combat.AttackCB) {
		trg := a.Target
		if trg.Type() != combat.TargettableEnemy {
			return
		}
		if !c.StatusIsActive(c4Key) {
			return
		}

		c.Core.Log.NewEvent("c4 proc'd on attack", glog.LogCharacterEvent, c.Index).
			Write("char", c.Index)
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Beidou C4",
			AttackTag:  attacks.AttackTagNone,
			ICDTag:     combat.ICDTagElementalBurst,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       0.2,
		}
		c.Core.QueueAttack(ai, combat.NewSingleTargetHit(trg.Key()), 0, 1)
	}
}
