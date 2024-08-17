package beidou

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

const c4Key = "beidou-c4"

// Upon being attacked, Beidou's Normal Attacks gain an additional instance of 20% Electro DMG for 10s.
func (c *char) c4Init() {
	c.Core.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
		di := args[0].(*info.DrainInfo)
		if !di.External {
			return false
		}
		if c.Core.Player.Active() != c.Index {
			return false
		}
		c.AddStatus(c4Key, 10*60, true)
		c.Core.Log.NewEvent("c4 triggered on damage", glog.LogCharacterEvent, c.Index).
			Write("expiry", c.StatusExpiry(c4Key))
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
		if trg.Type() != targets.TargettableEnemy {
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
			ICDTag:     attacks.ICDTagElementalBurst,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       0.2,
		}
		c.Core.QueueAttack(ai, combat.NewSingleTargetHit(trg.Key()), 0, 1)
	}
}
