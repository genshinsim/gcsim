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
func (c *char) c2(src int) func() {
	return func() {
		// don't tick if src changed
		if c.qFieldSrc != src {
			c.Core.Log.NewEvent("kazuha q src check ignored, src diff", glog.LogCharacterEvent, c.Index).
				Write("src", src).
				Write("new src", c.qFieldSrc)
			return
		}
		// don't tick if Q isn't up anymore
		if c.Core.Status.Duration(burstStatus) == 0 {
			return
		}

		// check again in 0.5s
		c.Core.Tasks.Add(c.c2(src), 30)

		ap := combat.NewCircleHitOnTarget(c.qAbsorbCheckLocation.Shape.Pos(), nil, 9)
		if !combat.TargetIsWithinArea(c.Core.Combat.Player().Pos(), ap) {
			return
		}

		// apply buff if in burst area
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
	}
}

// C6
// After using Chihayaburu or Kazuha Slash, Kaedehara Kazuha gains an Anemo
// Infusion for 5s. Additionally, each point of Elemental Mastery will
// increase the DMG dealt by Kaedehara Kazuha's Normal, Charged, and Plunging
// Attacks by 0.2%.
func (c *char) c6() {
	// add anemo infusion
	c.Core.Player.AddWeaponInfuse(
		c.Index,
		"kazuha-c6-infusion",
		attributes.Anemo,
		60*5,
		true,
		combat.AttackTagNormal, combat.AttackTagExtra, combat.AttackTagPlunge,
	)
	// add em based buff
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag("kazuha-c6-dmgup", 60*5), // 5s
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			// skip if not normal/charged/plunge
			if atk.Info.AttackTag != combat.AttackTagNormal &&
				atk.Info.AttackTag != combat.AttackTagExtra &&
				atk.Info.AttackTag != combat.AttackTagPlunge {
				return nil, false
			}
			// apply buff
			m[attributes.DmgP] = 0.002 * c.Stat(attributes.EM)
			return m, true
		},
	})
}
