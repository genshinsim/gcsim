package odette

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	stellarBonusKey  = "odette-stellar-bonus"
	radianceSwirlKey = "radiance-stellar-swirl"

	stellarConductText = " (Stellar-Conduct)"
	stellarSwirlText   = " (Stellar Swirl)"
)

type radianceState int

const (
	radianceNone radianceState = iota
	radianceStellarConduct
	radianceStellarSwirl
)

func (c *char) getRadiance() radianceState {
	if c.StatusIsActive(reactable.PolestarFieldKey) {
		return radianceStellarConduct
	}

	if c.StatusIsActive(radianceSwirlKey) {
		return radianceStellarSwirl
	}

	return radianceNone
}

// Odette will enter the Radiance: Stellar-Conduct state when she is inside a Polestar Field, or the
// Radiance: Stellar Swirl state for 8s after a nearby party member triggers a Stellar Swirl
// reaction.
// When a party member triggers a Superconduct or Cryo Swirl reaction, it becomes a Stellar-Conduct
// or Stellar Swirl reaction instead, and the Base DMG of said reaction is also increased by 0.7%
// for every 100 points of Odette's ATK. A maximum increase of 14% can be obtained in this way.
func (c *char) stellarInit() {
	c.Core.Flags.Custom[reactable.StellarConductEnableKey] = 1
	c.Core.Flags.Custom[reactable.StellarSwirlEnableKey] = 1

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectStellarConduct:
		case attacks.AttackTagDirectStellarSwirl:
		default:
			return
		}

		bonus := min(c.TotalAtk()/100.0*0.007, 0.14)

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("odette adding stellar base damage", glog.LogCharacterEvent, c.Index()).Write("bonus", bonus)
		}

		atk.Info.BaseDmgBonus += bonus
	}, stellarBonusKey)

	c.Core.Events.Subscribe(event.OnSpecialReactionAttack, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.AttackTag != attacks.AttackTagReactionStellarSwirl {
			return
		}

		bonus := min(c.TotalAtk()/100.0*0.007, 0.14)

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("odette adding stellar base damage", glog.LogCharacterEvent, c.Index()).Write("bonus", bonus)
		}

		atk.Info.BaseDmgBonus += bonus
	}, stellarBonusKey+"-reaction")

	c.Core.Events.Subscribe(event.OnStellarSwirl, func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		c.AddStatus(radianceSwirlKey, 8*60, false)
	}, stellarBonusKey)
}
