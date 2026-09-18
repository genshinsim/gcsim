package cryo

import (
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

type radianceState int

const (
	radianceNone radianceState = iota
	radianceStellarConduct
	radianceStellarSwirl
)

const (
	stellarBonusKey  = "travelercryo-stellar-bonus"
	radianceSwirlKey = "radiance-stellar-swirl"
)

func (c *Traveler) getRadiance() radianceState {
	if c.StatusIsActive(reactable.PolestarFieldKey) {
		return radianceStellarConduct
	}

	if c.StatusIsActive(radianceSwirlKey) {
		return radianceStellarSwirl
	}

	return radianceNone
}

// The Traveler will enter the Radiance: Stellar-Conduct state when inside a
// Polestar Field, or the Radiance: Stellar Swirl state for 8s after a nearby
// party member triggers a Stellar Swirl reaction.
//
// When a party member triggers a Superconduct or Cryo Swirl reaction, it
// becomes a Stellar-Conduct or Stellar Swirl reaction instead, and the Base DMG
// of said reaction is also increased by 0.35% for every 100 points of the
// Traveler's ATK. A maximum 7% increase can be obtained in this way.
func (c *Traveler) stellarInit() {
	c.Core.Flags.Custom[reactable.StellarConductEnableKey] = 1
	c.Core.Flags.Custom[reactable.StellarSwirlEnableKey] = 1

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		if !atk.Info.AttackTag.IsStellarDirect() {
			return
		}

		bonus := min(c.TotalAtk()/100.0*0.0035, 0.07)

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("travelercryo adding stellar base damage", glog.LogCharacterEvent, c.Index()).Write("bonus", bonus)
		}

		atk.Info.BaseDmgBonus += bonus
	}, stellarBonusKey)

	c.Core.Events.Subscribe(event.OnSpecialReactionAttack, func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		if !atk.Info.AttackTag.IsStellarReact() {
			return
		}

		bonus := min(c.TotalAtk()/100.0*0.0035, 0.07)

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("travelercryo adding stellar base damage", glog.LogCharacterEvent, c.Index()).Write("bonus", bonus)
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
