package sandrone

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

type radianceState int

const (
	radianceNone radianceState = iota
	radianceStellarConduct
	radianceStellarSwirl
)

const (
	stellarConductText = " (Stellar-Conduct)"
	stellarSwirlText   = " (Stellar-Swirl)"
)

const stellarBonusKey = "sandrone-stellar-bonus"

func (c *char) getRadiance() radianceState {
	if c.StatusIsActive(reactable.PolestarFieldKey) {
		return radianceStellarConduct
	}

	if c.StatusIsActive(reactable.RadianceSwirlKey) {
		return radianceStellarSwirl
	}

	return radianceNone
}

func (c *char) toStellar(ai *info.AttackInfo, sscMult, sswMult float64) {
	switch c.getRadiance() {
	case radianceStellarConduct:
		ai.Abil += stellarConductText
		ai.AttackTag = attacks.AttackTagDirectStellarConduct
		ai.Durability = 0
		ai.Mult = sscMult
		ai.IgnoreDefPercent = 1
	case radianceStellarSwirl:
		ai.Abil += stellarSwirlText
		ai.AttackTag = attacks.AttackTagDirectStellarSwirl
		ai.Durability = 0
		ai.Mult = sswMult
		ai.IgnoreDefPercent = 1
	}
}

// When a party member triggers a Superconduct or Cryo Swirl reaction, it becomes a Stellar-Conduct
// or Stellar Swirl reaction instead, and the Base DMG of the aforementioned reaction is also
// increased by 0.7% for every 100 points of Sandrone's ATK. A maximum increase of 14% can be
// obtained in this way.
func (c *char) stellarInit() {
	c.Core.Flags.Custom[reactable.StellarConductEnableKey] = 1
	c.Core.Flags.Custom[reactable.StellarSwirlEnableKey] = 1

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		if !atk.Info.AttackTag.IsStellarDirect() {
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

		if !atk.Info.AttackTag.IsStellarReact() {
			return
		}

		bonus := min(c.TotalAtk()/100.0*0.007, 0.14)

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("sandrone adding stellar base damage", glog.LogCharacterEvent, c.Index()).Write("bonus", bonus)
		}

		atk.Info.BaseDmgBonus += bonus
	}, stellarBonusKey+"-reaction")
}
