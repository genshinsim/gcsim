package sandrone

import (
	"slices"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1Key    = "sandrone-c1"
	c2Key    = "sandrone-c2"
	c4Key    = "sandrone-c4"
	c4ICDKey = "sandrone-c4-icd"
	c6Key    = "sandrone-c6"

	// TODO: Get c6 hitmark
	c6Hitmark = 3
)

// the first beam already has a stack so the 0.4 is never used
var c2Buff = []float64{0.4, 0.6, 0.8, 1.0}

// When in Decoding mode, Fagio's Decoding Power increases at a 50% lower rate, and all party
// members deal 30% increased Stellar Glimmer reaction DMG.
func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		return
	}

	c.caMaxBeams = 6
	for _, char := range c.Core.Player.Chars() {
		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase(c1Key, -1),
			Amount: func(ai info.AttackInfo) float64 {
				if c.currFagio != stateDecoding {
					return 0
				}

				if !ai.AttackTag.IsStellar() {
					return 0
				}

				return 0.3
			},
		})
	}
}

func (c *char) c1DecoderGainMult() float64 {
	if c.Base.Cons < 1 {
		return 1.0
	}
	return 0.5
}

// Radiance: Stellar Glimmer: Increases CRIT DMG from Charged Attack condensed beams by 40%. Also,
// each time a beam is fired, CRIT DMG dealt by all condensed beams fired during this Decoding mode
// period will be further increased by 20%. Max 3 stacks.
func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(c2Key, -1),
		Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
			if !slices.Contains(atk.Info.AdditionalTags, attacks.AttackTagSandroneBeam) {
				return nil
			}
			if c.getRadiance() == radianceNone {
				return nil
			}
			m[attributes.CD] = c2Buff[c.c2stacks]
			return m
		},
	})
}

// stack gain before beam hits
func (c *char) c2OnBeam() {
	if c.Base.Cons < 2 {
		return
	}

	c.c2stacks = min(c.c2stacks+1, 3)
}

func (c *char) c2OnCAStart() {
	if c.Base.Cons < 2 {
		return
	}
	c.c2stacks = 0
}

// When Sandrone deals Stellar-Conduct or Stellar Swirl DMG to an opponent, she will also summon an
// extra Prismatic Resonance Cannon to fire off a coordinated attack which deals Cryo DMG at 125% or
// 187.5% of her ATK, respectively. This DMG is considered the corresponding Stellar Glimmer
// reaction DMG, and the effect can occur once every 4s.
func (c *char) c4Init() {
	if c.Base.Cons < 4 {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		e, ok := args[0].(*enemy.Enemy)
		if !ok {
			return
		}
		atk, ok := args[1].(*info.AttackEvent)
		if !ok {
			return
		}
		if atk.Info.ActorIndex != c.Index() {
			return
		}

		var mult float64
		var tag attacks.AttackTag
		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectStellarConduct:
			mult = 1.25
			tag = attacks.AttackTagDirectStellarConduct
		case attacks.AttackTagDirectStellarSwirl, attacks.AttackTagReactionStellarSwirl:
			mult = 1.875
			tag = attacks.AttackTagDirectStellarSwirl
		default:
			return
		}

		if c.StatusIsActive(c4ICDKey) {
			return
		}

		c.AddStatus(c4ICDKey, 4*60, true)

		ai := info.AttackInfo{
			ActorIndex:       c.Index(),
			Abil:             "Prismatic Resonance Cannon (C4)",
			AttackTag:        tag,
			Element:          attributes.Cryo,
			Mult:             mult,
			IgnoreDefPercent: 1,
		}

		ap := combat.NewSingleTargetHit(e.Key())

		// TODO: Get c4 hitmark
		c.Core.QueueAttack(ai, ap, 15, 15)
	}, c4Key)
}

// When Fagio fires a condensed beam for the third time while in Decoding mode, the beam becomes a
// more powerful Condensed Cluster Beam that fires continuously. This beam will deal up to 4
// additional instances of AoE Cryo DMG at 100% of Sandrone's ATK on top of any subsequent original
// beam fired.
// Radiance: Stellar Glimmer: The Condensed Cluster Beam attack deals 4 additional instances of AoE
// Cryo DMG. The DMG dealt with this attack, which is considered the corresponding Stellar Glimmer
// reaction DMG, is equal to:
// Radiance: Stellar-Conduct: 80% of Sandrone's ATK.
// Radiance: Stellar Swirl: 120% of Sandrone's ATK.
// Additionally, all Stellar Glimmer reaction DMG dealt by Sandrone is elevated by 20%.
func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}

	c.Core.Events.Subscribe(event.OnApplyAttack, func(args ...any) {
		atk := args[0].(*info.AttackEvent)

		// don't apply elevation to the reaction attack, since the subcomponent contributor attacks each got elevation applied already
		if !atk.Info.AttackTag.IsStellarDirect() {
			return
		}

		if atk.Info.ActorIndex != c.Index() {
			return
		}

		if c.getRadiance() == radianceNone {
			return
		}

		atk.Info.Elevation += 0.2
	}, c6Key)

	c.Core.Events.Subscribe(event.OnSpecialReactionAttack, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if !atk.Info.AttackTag.IsStellarReact() {
			return
		}

		if atk.Info.ActorIndex != c.Index() {
			return
		}

		if c.getRadiance() == radianceNone {
			return
		}

		atk.Info.Elevation += 0.2
	}, c6Key+"-stellar-atk")
}

func (c *char) c6OnCaStart() {
	if c.Base.Cons < 6 {
		return
	}
	c.c6stacks = 0
}

// stack gain before beam hits
func (c *char) c6OnBeam() {
	if c.Base.Cons < 6 {
		return
	}

	c.c6stacks += 1
	if c.c6stacks < 3 {
		return
	}

	ai := info.AttackInfo{
		ActorIndex:     c.Index(),
		Abil:           "Condensed Cluster Beam (C6)",
		AttackTag:      attacks.AttackTagExtra,
		AdditionalTags: []attacks.AttackTag{attacks.AttackTagSandroneBeam},
		ICDTag:         attacks.ICDTagSandroneExtraAttackLaser,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeBlunt,
		PoiseDMG:       45,
		Element:        attributes.Cryo,
		Durability:     25,
		Mult:           1,
	}

	c.toStellar(&ai, 0.8, 1.2)

	ap := combat.NewBoxHitOnTarget(
		c.Core.Combat.Player(),
		info.Point{Y: -7.5},
		3,
		15,
	)

	c.Core.QueueAttack(ai, ap, c6Hitmark, c6Hitmark)
}
