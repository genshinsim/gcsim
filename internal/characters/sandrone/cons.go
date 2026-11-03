package sandrone

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1Key      = "sandrone-c1"
	c2Key      = "sandrone-c2"
	c2MaxStack = 3
	c4Key      = "sandrone-c4"
	c4ICDKey   = "sandrone-c4-icd"
	c4Abil     = "Coordinated Attack (C4)"
	c6Key      = "sandrone-c6"
	c6Abil     = "Condensed Cluster Beam (C6)"
	c6Hits     = 4
	c6HitGap   = 6
	c6MinBeams = 3
)

// When in Decoding mode, Fagio's Decoding Power increases at a 50% lower rate (see decodeMult),
// and all party members deal 30% increased Stellar Glimmer reaction DMG.
func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		return
	}
	for _, char := range c.Core.Player.Chars() {
		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase(c1Key, -1),
			Amount: func(ai info.AttackInfo) float64 {
				if c.fagioMode != fagioDecoding {
					return 0
				}
				if !ai.AttackTag.IsStellar() {
					return 0
				}
				return c1Bonus
			},
		})
	}
}

// Radiance: Stellar Glimmer: Increases CRIT DMG from Charged Attack condensed beams by 40%. Also,
// each time a beam is fired, CRIT DMG dealt by all condensed beams fired during this Decoding mode
// period will be further increased by 20%. Max 3 stacks.
// Stacks are counted from the beams fired before the current one in this Decoding period.
func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(c2Key, -1),
		Amount: func(atk *info.AttackEvent, _ info.Target) []float64 {
			if !atk.Info.AttackTag.IsStellarDirect() {
				return nil
			}
			if !strings.HasPrefix(atk.Info.Abil, beamAbil) {
				return nil
			}
			stacks := min(max(c.beamCount-1, 0), c2MaxStack)
			m[attributes.CD] = c2CD + c2CDStack*float64(stacks)
			return m
		},
	})
}

// When Sandrone deals Stellar-Conduct or Stellar Swirl DMG to an opponent, she will also summon an
// extra Prismatic Resonance Cannon to fire off a coordinated attack which deals Cryo DMG at 125% or
// 187.5% of her ATK, respectively. This DMG is considered the corresponding Stellar Glimmer
// reaction DMG, and the effect can occur once every 4s.
func (c *char) c4Init() {
	if c.Base.Cons < 4 {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		trg, ok := args[0].(info.Target)
		if !ok {
			return
		}
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != c.Index() || atk.Info.Abil == c4Abil {
			return
		}
		mult := 0.0
		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectStellarConduct:
			mult = c4Conduct
		case attacks.AttackTagDirectStellarSwirl:
			mult = c4Swirl
		default:
			return
		}
		if c.StatusIsActive(c4ICDKey) {
			return
		}
		c.AddStatus(c4ICDKey, c4CD*60, true)

		ai := info.AttackInfo{
			ActorIndex:       c.Index(),
			Abil:             c4Abil,
			AttackTag:        atk.Info.AttackTag,
			ICDTag:           attacks.ICDTagNone,
			ICDGroup:         attacks.ICDGroupDefault,
			StrikeType:       attacks.StrikeTypeDefault,
			Element:          attributes.Cryo,
			Mult:             mult,
			IgnoreDefPercent: 1,
		}
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(trg, nil, 2), 0, 0)
	}, c4Key)
}

// When Fagio fires a condensed beam for the third time while in Decoding mode, the beam becomes a
// more powerful Condensed Cluster Beam: up to 4 additional instances of AoE Cryo DMG at 100% of
// Sandrone's ATK on top of any subsequent original beam fired.
// Radiance: Stellar Glimmer: these instances are Stellar Glimmer reaction DMG at 80% (Stellar-
// Conduct) or 120% (Stellar Swirl) of Sandrone's ATK.
func (c *char) c6ClusterBeam(ap info.AttackPattern) {
	if c.Base.Cons < 6 || c.beamCount < c6MinBeams {
		return
	}
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       c6Abil,
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagExtraAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       c6Cluster,
	}
	c.toStellar(&ai, c6Conduct, c6Swirl)
	for i := 1; i <= c6Hits; i++ {
		c.Core.QueueAttack(ai, ap, 0, i*c6HitGap)
	}
}

// Additionally, all Stellar Glimmer reaction DMG dealt by Sandrone is elevated by 20%.
func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}
	c.Core.Events.Subscribe(event.OnApplyAttack, func(args ...any) {
		atk := args[0].(*info.AttackEvent)
		if atk.Info.ActorIndex != c.Index() {
			return
		}
		if atk.Info.AttackTag.IsStellarDirect() {
			atk.Info.Elevation += c6Elevate
		}
	}, c6Key+"-direct")

	c.Core.Events.Subscribe(event.OnSpecialReactionAttack, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != c.Index() {
			return
		}
		if atk.Info.AttackTag.IsStellarReact() {
			atk.Info.Elevation += c6Elevate
		}
	}, c6Key+"-reaction")
}
