package kazuha

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// C2:
// The Autumn Whirlwind field created by Kazuha Slash has the following effects:
// - Increases Kaedehara Kazuha's own Elemental Mastery by 200 for its duration.
// - Increases the Elemental Mastery of characters within the field by 200.
// The Elemental Mastery-increasing effects of this Constellation do not stack.
func (c *char) c2() {
	// don't tick if Q isn't up anymore
	if c.Core.Status.Duration(burstStatus) == 0 {
		return
	}

	c.Core.Log.NewEvent("kazuha-c2 ticking", glog.LogCharacterEvent, -1)

	// apply C2 buff to active char for 1s
	active := c.Core.Player.ActiveChar()
	active.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("kazuha-c2", 60), // 1s
		AffectedStat: attributes.EM,
		Amount: func() ([]float64, bool) {
			return c.c2buff, true
		},
	})

	// apply C2 buff to Kazuha (even if off-field) for 1s
	if active.Base.Key != c.Base.Key {
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("kazuha-c2", 60), // 1s
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return c.c2buff, true
			},
		})
	}

	// check again in 0.5s
	c.Core.Tasks.Add(c.c2, 30)
}

// C6
// After using Chihayaburu or Kazuha Slash, Kaedehara Kazuha gains an Anemo
// Infusion for 5s. Additionally, each point of Elemental Mastery will
// increase the DMG dealt by Kaedehara Kazuha's Normal, Charged, and Plunging
// Attacks by 0.2%.
func (c *char) c6() {
	c.AddStatus(c6BuffKey, 60*5, true)
	c.Core.Player.AddWeaponInfuse(
		c.Index,
		"kazuha-c6-infusion",
		attributes.Anemo,
		60*5,
		true,
		combat.AttackTagNormal, combat.AttackTagExtra, combat.AttackTagPlunge,
	)
}
