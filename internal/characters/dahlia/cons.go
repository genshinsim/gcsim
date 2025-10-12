package dahlia

import (
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

const (
	c1Key = "dahlia-c1"
	c2Key = "dahlia-c2"
	c6Key = "dahlia-c6"
	c6Icd = "dahlia-c6-icd"
)

// C1
// Each time Dahlia gains 1 of his Elemental Burst Radiant Psalter's Benison stacks, he will regain 2.5 Elemental Energy.
func (c *char) c1OnBenisonEnergy(stacks int) {
	if c.Base.Cons < 1 {
		return
	}
	c.AddEnergy(c1Key, 2.5*float64(stacks))
}

// C2
// After Dahlia consumes his Elemental Burst Radiant Psalter's Benison stacks to summon a Shield of Sacred Favor,
// the character protected by said Shield will gain 25% increased Shield Strength for 12s.
// NOTE: Mechanics in game
// - Shield strength gets refreshed if C2 is triggered again
// - Shield strength is given to whoever is currently on-field (even if C2 triggered when another character was on-field)
// - Shield strength is only applied while shielded by Dahlia
// - Shield strength carries over to the initial shield from a subsequent Burst
func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}

	c.Core.Player.Shields.AddShieldBonusMod(c2Key, 12*60, func() (float64, bool) {
		if c.hasShield() { // TO-DO: does this do what the 3rd bullet point is saying?
			return 0.25, false
		}
		return 0, false
	})
}

// C4
// The Favonian Favor from Dahlia's Elemental Burst Radiant Psalter lasts 3 more seconds.
func (c *char) c4FavonianFavorBonusDur() int {
	if c.Base.Cons < 4 {
		return 0
	}
	return 3 * 60
}

// C6
// The current active character affected by the Elemental Burst Radiant Psalter's Favonian Favor has their
// ATK SPD increased by 10%.
// Additionally, when an active character affected by Favonian Favor falls, immediately:
// - Revive them.
// - Restore their HP to 100%.
// This effect can trigger once every 15 minutes.
func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}

	c.Core.Events.Subscribe(event.OnPlayerHPDrain, func(args ...any) bool {
		char := c.Core.Player.ActiveChar()

		if !char.StatusIsActive(burstFavonianFavor) {
			return false
		}

		di := args[0].(*info.DrainInfo)
		if di.Amount <= 0 {
			return false
		}

		// If Revive is still on CD, do nothing
		if c.StatusIsActive(c6Icd) {
			return false
		}

		// Revive the active char (even Dahlia himself) back to 100% HP if dead
		if char.CurrentHPRatio() <= 0 {
			char.SetHPByRatio(1)
		}
		c.AddStatus(c6Icd, 15*60*60, true)

		return false
	}, c6Key)
}

func (c *char) c6AtkSpd() float64 {
	if c.Base.Cons < 6 {
		return 0.0
	}
	return 0.1
}
