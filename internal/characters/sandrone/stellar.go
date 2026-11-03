package sandrone

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	stellarBonusKey = "sandrone-stellar-bonus"

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

	if c.StatusIsActive(reactable.RadianceSwirlKey) {
		return radianceStellarSwirl
	}

	return radianceNone
}

// toStellar converts an attack into the Stellar Glimmer reaction DMG variant matching the current
// Radiance state. Returns false (and leaves ai untouched) when Sandrone is not in Radiance.
func (c *char) toStellar(ai *info.AttackInfo, conduct, swirl float64) bool {
	switch c.getRadiance() {
	case radianceStellarConduct:
		ai.Abil += stellarConductText
		ai.AttackTag = attacks.AttackTagDirectStellarConduct
		ai.Mult = conduct
	case radianceStellarSwirl:
		ai.Abil += stellarSwirlText
		ai.AttackTag = attacks.AttackTagDirectStellarSwirl
		ai.Mult = swirl
	default:
		return false
	}
	ai.ICDTag = attacks.ICDTagNone
	ai.ICDGroup = attacks.ICDGroupDefault
	ai.Durability = 0
	ai.IgnoreDefPercent = 1
	return true
}

// Sandrone will enter the Radiance: Stellar-Conduct state when she is inside a Polestar Field, or
// the Radiance: Stellar Swirl state for 8s after a nearby party member triggers a Stellar Swirl
// reaction.
// When a party member triggers a Superconduct or Cryo Swirl reaction, it becomes a Stellar-Conduct
// or Stellar Swirl reaction instead, and the Base DMG of the aforementioned reaction is also
// increased by 0.7% for every 100 points of Sandrone's ATK. A maximum increase of 14% can be
// obtained in this way.
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
			c.Core.Log.NewEvent("sandrone adding stellar base damage", glog.LogCharacterEvent, c.Index()).Write("bonus", bonus)
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
			c.Core.Log.NewEvent("sandrone adding stellar base damage", glog.LogCharacterEvent, c.Index()).Write("bonus", bonus)
		}

		atk.Info.BaseDmgBonus += bonus
	}, stellarBonusKey+"-reaction")
}
