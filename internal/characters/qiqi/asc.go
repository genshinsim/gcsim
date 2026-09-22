package qiqi

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const a4ICDKey = "qiqi-a4-icd"

// When a character under the effects of Adeptus Art: Herald of Frost triggers an Elemental Reaction,
// their Incoming Healing Bonus is increased by 20% for 8s.
// - implements event hook and incoming healing bonus function
// - TODO: Could possibly change this so the AddIncHealBonus occurs at start, then event subscription occurs upon using Qiqi skill?
// - TODO: Likely more efficient to not maintain event subscription always, but grouping the two for clarity currently
func (c *char) a1Init() {
	if c.Base.Ascension < 1 {
		return
	}
	a1Hook := func(args ...any) {
		if !c.StatusIsActive(skillBuffKey) {
			return
		}
		atk := args[1].(*info.AttackEvent)

		// Active char is the only one under the effects of Qiqi skill
		active := c.Core.Player.ActiveChar()
		if atk.Info.ActorIndex != active.Index() {
			return
		}

		active.AddHealBonusMod(character.HealBonusMod{
			Base: modifier.NewBaseWithHitlag("qiqi-a1", 8*60),
			Amount: func() float64 {
				return .2
			},
		})
	}

	for i := event.ReactionEventStartDelim + 1; i < event.ReactionEventEndDelim; i++ {
		c.Core.Events.Subscribe(i, a1Hook, "qiqi-a1")
	}
}

func (c *char) revelationA4() (float64, int) {
	if !c.revelation {
		return 0.5, 30 * 60
	}

	if c.getRadiance() == radianceNone {
		return 0.5, 30 * 60
	}

	return 1.0, 15 * 60
}

// A4 is implemented in burst.go:
// When Qiqi hits opponents with her Normal and Charged Attacks,
// she has a 50% chance to apply a Fortune-Preserving Talisman to them for 6s.
// This effect can only occur once every 30s.
//
// Radiance: Stellar Glimmer: The effect now has a 100% chance to trigger,
// and cooldown time is reduced to 15s.
func (c *char) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		e, ok := args[0].(*enemy.Enemy)
		atk := args[1].(*info.AttackEvent)
		if !ok {
			return
		}
		if atk.Info.ActorIndex != c.Index() {
			return
		}

		// All of the below only occur on Qiqi NA/CA hits
		switch atk.Info.AttackTag {
		case attacks.AttackTagNormal:
		case attacks.AttackTagExtra:
		default:
			return
		}

		if c.StatusIsActive(a4ICDKey) {
			return
		}

		chance, cd := c.revelationA4()

		if c.Core.Rand.Float64() > chance {
			return
		}

		// Don't want to overwrite a longer burst duration talisman with a shorter duration one
		// TODO: Unclear how the interaction works if there is already a talisman on enemy
		// TODO: Being generous for now and not putting it on CD if there is a conflict
		if e.StatusExpiry(talismanKey) < c.Core.F+360 {
			e.AddStatus(talismanKey, 360, true)
			c.AddStatus(a4ICDKey, cd, true)
			c.Core.Log.NewEvent(
				"Qiqi A4 Adding Talisman",
				glog.LogCharacterEvent,
				c.Index(),
			).
				Write("target", e.Key()).
				Write("talisman_expiry", e.StatusExpiry(talismanKey))
		}
	}, "qiqi-a4")
}
