package vodyanitsa

import (
	"math"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	a1Key = "vodyanitsa-a1"
	a4Key = "vodyanitsa-a4"
)

// When the Elemental Skill Rechitativ: Sonorous Dawn is used, any Stellar Vortex on the field will
// turn into a Wandering Vortex instead. While the Song of Ages Past effect is active, nearby party
// members' Stellar Swirl reactions will generate Wandering Vortex instead.
// Generating or detonating a Wandering Vortex will also reduce nearby opponents' Anemo RES by 35%
// for 6s.
// Additionally, nearby party members who enter the Radiance: Stellar Swirl state while Song of Ages
// Past is in effect will have the Radiance: Stellar Swirl state extended by 4s.
func (c *char) a1Init() {
	if c.Base.Ascension < 1 {
		return
	}
	c.nearbyWanderingVortexExp = -1

	c.Core.Events.Subscribe(event.OnStellarVortex, func(args ...any) {
		if !c.StatusIsActive(skillKey) {
			return
		}
		c.nearbyWanderingVortexExp = c.Core.F + math.MaxUint16
		c.a1ApplyShred()
	}, a1Key)

	c.Core.Events.Subscribe(event.OnStellarVortexDetonate, func(args ...any) {
		if !c.StatusIsActive(skillKey) {
			return
		}
		c.nearbyWanderingVortexExp = c.Core.F + 5*60
		c.a1ApplyShred()
	}, "vodyanitsa-ssw")
}

func (c *char) recentWandering() bool {
	return c.Core.F < c.nearbyWanderingVortexExp
}

func (c *char) a1ApplyShred() {
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5)
	for _, e := range c.Core.Combat.EnemiesWithinArea(ap, nil) {
		e, ok := e.(*enemy.Enemy)
		if !ok {
			continue
		}
		e.AddResistMod(info.ResistMod{
			Base:  modifier.NewBaseWithHitlag(a1Key, 6*60),
			Ele:   attributes.Anemo,
			Value: -0.35,
		})
	}
}

func (c *char) a1OnSkill(src int) {
	if c.Base.Ascension < 1 {
		return
	}

	for _, g := range c.Core.Combat.Gadgets() {
		if g != nil && g.GadgetTyp() == info.GadgetTypStellarVortex {
			c.a1ApplyShred()
			break
		}
	}

	c.Core.Flags.Custom[reactable.RadianceSwirlBonusDurKey] = 4 * 60
	c.Core.Tasks.Add(func() {
		if c.skillSrc != src {
			return
		}
		// assumes no one else uses this flag
		c.Core.Flags.Custom[reactable.RadianceSwirlBonusDurKey] = 0
	}, skillDur+c.c2SkillDurBonus())
}

// When the Elemental Skill Rechitativ: Sonorous Dawn is used, Vodyanitsa will also gain 25 stacks
// of Lead Vocal as well as 10 stacks of Chorus.
// The two effects above last 30s each, and the stack count is refreshed when the Elemental Skill is
// used.
//
// Lead Vocal:
// When the party's currently active character deals Hydro DMG or Cryo DMG with a Normal Attack,
// Charged Attack, Plunging Attack, Elemental Skill, or Elemental Burst, 1 stack of "Lead Vocal" is
// consumed. This in turn increases the DMG dealt. However, if a Wandering Vortex is present on the
// field or during the 5s following a Wandering Vortex explosion, this effect will only trigger when
// the character deals Stellar Swirl reaction DMG.
// Every 1000 Max HP Vodyanitsa has in excess of 40000 will increase the Stellar Swirl reaction DMG
// from the above effect by 260, as well as increase Hydro DMG and Cryo DMG dealt by 140. Stellar
// Swirl reaction DMG dealt can be increased by a maximum of 6500, while Hydro DMG and Cryo DMG can
// be increased by a maximum of 3500 in this way.
//
// Chorus:
// Triggered by off-field party members or other nearby party members (and not your currently active
// party member). Effects are the same as those of "Lead Vocal."
func (c *char) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}

	activeContrib := false
	offFieldContrib := false

	// this is emitted up to 4 times per attack, but we should only consume 1 stack per "attack"
	// so we only add damage in OnLunarReactionAttack, but we will consume in OnEnemyHit
	// This relies on OnSpecialReactionAttack to be emitted before the corresponding OnEnemyHit,
	// even in AoE.
	c.Core.Events.Subscribe(event.OnSpecialReactionAttack, func(args ...any) {
		ae := args[1].(*info.AttackEvent)
		owner := ae.Info.ActorIndex
		if !c.StatusIsActive(a4Key) {
			return
		}

		var amt float64
		if ae.Info.AttackTag != attacks.AttackTagReactionStellarSwirl {
			return
		}

		if !c.a4StacksCheck(owner) {
			return
		}

		amt = max(min((c.MaxHP()-40000)*260/1000, 6500), 0)
		ae.Info.FlatDmg += amt

		if owner == c.Core.Player.Active() {
			activeContrib = true
		} else {
			offFieldContrib = true
		}
	}, a4Key)

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		ae := args[1].(*info.AttackEvent)

		if !c.StatusIsActive(a4Key) {
			return
		}

		var amt float64
		switch ae.Info.AttackTag {
		case attacks.AttackTagNormal,
			attacks.AttackTagExtra, attacks.AttackTagPlunge,
			attacks.AttackTagElementalArt,
			attacks.AttackTagElementalArtHold,
			attacks.AttackTagElementalBurst:
			if ae.Info.Element != attributes.Cryo && ae.Info.Element != attributes.Hydro {
				return
			}
			if c.recentWandering() {
				return
			}
			if !c.a4StacksCheck(ae.Info.ActorIndex) {
				return
			}
			amt = max(min((c.MaxHP()-40000)*140/1000, 3500), 0)
		case attacks.AttackTagDirectStellarSwirl:
			if !c.recentWandering() {
				return
			}
			if !c.a4StacksCheck(ae.Info.ActorIndex) {
				return
			}
			amt = max(min((c.MaxHP()-40000)*260/1000, 6500), 0)
		case attacks.AttackTagReactionStellarSwirl:
			if !c.recentWandering() {
				return
			}

			// Based on testing, even if there are multiple off field contributors, only 1 stack
			// of Chorus is consumed per reaction attack
			c.a4StacksRemoveReaction(activeContrib, offFieldContrib)
			activeContrib = false
			offFieldContrib = false
			return
		default:
			return
		}

		c.a4StacksRemove(ae.Info.ActorIndex)
		ae.Info.FlatDmg += amt
	}, a4Key)
}

func (c *char) a4StacksCheck(actor int) bool {
	if actor == c.Core.Player.Active() {
		return c.leadVocal > 0
	}
	return c.chorus > 0
}

func (c *char) a4StacksRemove(actor int) {
	if actor == c.Core.Player.Active() {
		c.leadVocal = max(c.leadVocal-1, 0)
	} else {
		c.chorus = max(c.chorus-1, 0)
	}
}

func (c *char) a4StacksRemoveReaction(lead, chorus bool) {
	if lead {
		c.leadVocal = max(c.leadVocal-1, 0)
	}
	if chorus {
		c.chorus = max(c.chorus-1, 0)
	}
}

func (c *char) a4OnSkill() {
	if c.Base.Ascension < 4 {
		return
	}

	c.leadVocal = 25
	c.chorus = 10
	c.AddStatus(a4Key, 30*60, true)
}
