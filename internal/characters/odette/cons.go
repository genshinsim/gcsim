package odette

import (
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
	c2Key    = "odette-c2"
	c4Key    = "odette-c4"
	c4ICDKey = "odette-c4-icd"
	c6Key    = "odette-c6"
)

// After unleashing the special Elemental Skill Adagio: Coda at Dawn's Tolling, at the dance duet's
// end Odette will deal an additional instance of Cryo AoE DMG to nearby opponents that is
// considered:
// - Radiance: Stellar-Conduct or when not in a Radiance state: Stellar-Conduct reaction DMG at 300% of Odette's ATK;
// - Radiance: Stellar Swirl: Stellar Swirl reaction DMG at 450% of Odette's ATK.

// Additionally, the Ascension Talent "Spring Rite of the Chosen One" is also enhanced: now, when
// the Solo Dance Double is summoned, Odette also gains 2 stacks of Marvelous Splendor. When Odette
// is off-field, the rate at which Marvelous Splendor is removed is sped up to 2 stacks per second.
func (c *char) c1OnSkillRecast(tag attacks.AttackTag) {
	if c.Base.Cons < 1 {
		return
	}

	ai := info.AttackInfo{
		ActorIndex:       c.Index(),
		AttackTag:        tag,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Cryo,
		IgnoreDefPercent: 1,
	}

	switch tag {
	case attacks.AttackTagDirectStellarConduct:
		ai.Abil = "Daybreak Finale (C1)" + stellarConductText
		ai.Mult = 3 * c.a4StellarGlimmerMult()
	case attacks.AttackTagDirectStellarSwirl:
		ai.Abil = "Daybreak Finale (C1)" + stellarSwirlText
		ai.Mult = 4.5 * c.a4StellarGlimmerMult()
	}
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6)
	c.Core.QueueAttack(ai, ap, 0, 0, c.particleCB)
}

func (c *char) c1a1Stacks() int {
	if c.Base.Cons < 1 {
		return 0
	}
	return 2
}

func (c *char) c1a1Remove() int {
	if c.Base.Cons < 1 {
		return 1
	}
	return 2
}

// The Ascension Talent "Spring Rite of the Chosen One" is enhanced as follows: every stack of
// Marvelous Splendor active also increases the character's ATK by 7%.
// Additionally, if Odette is in the Radiance: Stellar Glimmer state when there is a Solo Dance
// Double on the field, opponents near the Dance Double will also have their corresponding Elemental
// RES lowered by 20%.
// - Radiance: Stellar-Conduct: Cryo and Electro.
// - Radiance: Stellar Swirl: Cryo and Anemo.
func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}

	if c.Base.Ascension < 1 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(c2Key, -1),
		AffectedStat: attributes.ATKP,
		Amount: func() []float64 {
			if !c.StatusIsActive(danceDoubleKey) {
				return nil
			}
			m[attributes.ATKP] = float64(c.a1StacksSelf) * 0.07
			return m
		},
	})

	for _, char := range c.Core.Player.Chars() {
		if char.Index() == c.Index() {
			continue
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase(c2Key, -1),
			AffectedStat: attributes.ATKP,
			Amount: func() []float64 {
				if !c.StatusIsActive(danceDoubleKey) {
					return nil
				}
				m[attributes.ATKP] = float64(c.a1StacksOthers) * 0.07
				return m
			},
		})
	}
}

func (c *char) c2OnDanceSummon() {
	if c.Base.Cons < 2 {
		return
	}

	if c.Base.Ascension < 1 {
		return
	}
	c.c2Src = c.Core.F
	c.c2Ticker(c.c2Src)
}

func (c *char) c2Ticker(src int) {
	if !c.StatusIsActive(danceDoubleKey) {
		return
	}

	if c.c2Src != src {
		return
	}

	c.Core.Tasks.Add(func() { c.c2Ticker(src) }, 0.3*60)

	var otherElem attributes.Element
	switch c.getRadiance() {
	case radianceStellarConduct:
		otherElem = attributes.Electro
	case radianceStellarSwirl:
		otherElem = attributes.Anemo
	default:
		return
	}

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 10)
	for _, e := range c.Core.Combat.EnemiesWithinArea(ap, nil) {
		e, ok := e.(*enemy.Enemy)
		if !ok {
			continue
		}
		e.AddResistMod(info.ResistMod{
			Base:  modifier.NewBaseWithHitlag(c2Key+"-"+attributes.Cryo.String(), 1*60),
			Ele:   attributes.Cryo,
			Value: -0.20,
		})

		e.AddResistMod(info.ResistMod{
			Base:  modifier.NewBaseWithHitlag(c2Key+"-"+otherElem.String(), 1*60),
			Ele:   otherElem,
			Value: -0.20,
		})
	}
}

// The Elemental Burst Presto: Bluebird Finale is enhanced as follows: when Odette obtains Snow
// Swan's Dream, Stellar Glimmer reaction DMG dealt by other nearby party members is increased by
// 50% of Snow Swan Dream's effects.
func (c *char) c4OnBurst(buff float64) {
	if c.Base.Cons < 4 {
		return
	}
	buff *= 0.5
	for _, char := range c.Core.Player.Chars() {
		if char.Index() == c.Index() {
			continue
		}
		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBaseWithHitlag(swansDreamKey, 20*60),
			Amount: func(ai info.AttackInfo) float64 {
				switch ai.AttackTag {
				case
					attacks.AttackTagDirectStellarConduct,
					attacks.AttackTagDirectStellarSwirl,
					attacks.AttackTagReactionStellarSwirl:
					return buff
				default:
					return 0
				}
			},
		})
	}
}

// Additionally, when a party member deals Stellar Glimmer reaction DMG to an opponent, Odette will also join in with a coordinated attack, dealing an instance of AoE Cryo DMG. This effect, which can trigger once every 3.5s, will be considered:
// - Radiance: Stellar-Conduct or when not in a Radiance state: Stellar-Conduct reaction DMG at 66% of Odette's ATK.
// - Radiance: Stellar Swirl: Stellar Swirl reaction DMG at 99% of Odette's ATK.
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

		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectStellarConduct,
			attacks.AttackTagDirectStellarSwirl,
			attacks.AttackTagReactionStellarSwirl:
		default:
			return
		}

		if c.StatusIsActive(c4ICDKey) {
			return
		}

		c.AddStatus(c4ICDKey, 3.5*60, true)

		ai := info.AttackInfo{
			ActorIndex:       c.Index(),
			Abil:             "Odette C4",
			ICDTag:           attacks.ICDTagNone,
			ICDGroup:         attacks.ICDGroupDefault,
			StrikeType:       attacks.StrikeTypeDefault,
			Element:          attributes.Cryo,
			IgnoreDefPercent: 1,
		}

		switch c.getRadiance() {
		case radianceStellarConduct,
			radianceNone:
			ai.Abil += stellarConductText
			ai.AttackTag = attacks.AttackTagDirectStellarConduct
			ai.Mult = 0.66
		case radianceStellarSwirl:
			ai.Abil += stellarSwirlText
			ai.AttackTag = attacks.AttackTagDirectStellarSwirl
			ai.Mult = 0.99
		default:
			return
		}

		ap := combat.NewCircleHitOnTarget(e, nil, 2)
		c.Core.QueueAttack(ai, ap, 5, 5)
	}, c4Key)
}

// The Ascension Talent "Spring Rite of the Chosen One" is enhanced as follows: When Odette grants
// Marvelous Splendor to all nearby party members, her own Marvelous Splendor stacks will no longer
// decrease.
// Additionally, characters affected by Marvelous Splendor have their Stellar Glimmer reaction DMG
// dealt to opponents elevated by 25%, and Stellar Glimmer reaction DMG dealt by Odette is elevated
// by an additional 20%.
func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}

	odette := 0.45
	other := 0.25

	c.Core.Events.Subscribe(event.OnApplyAttack, func(args ...any) {
		atk := args[0].(*info.AttackEvent)

		// don't apply elevation to the reaction attack, since the subcomponent contributor attacks each got elevation applied already
		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectStellarConduct:
		case attacks.AttackTagDirectStellarSwirl:
		default:
			return
		}

		if !c.StatusIsActive(danceDoubleKey) {
			return
		}

		if atk.Info.ActorIndex == c.Index() {
			if c.a1StacksSelf > 0 {
				atk.Info.Elevation += odette
			}
			return
		}

		if c.a1StacksOthers > 0 {
			atk.Info.Elevation += other
		}
	}, c6Key)

	c.Core.Events.Subscribe(event.OnSpecialReactionAttack, func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		if atk.Info.AttackTag != attacks.AttackTagReactionStellarSwirl {
			return
		}

		if !c.StatusIsActive(danceDoubleKey) {
			return
		}

		if atk.Info.ActorIndex == c.Index() {
			if c.a1StacksSelf > 0 {
				atk.Info.Elevation += odette
			}
			return
		}

		if c.a1StacksOthers > 0 {
			atk.Info.Elevation += other
		}
	}, c6Key)
}

func (c *char) c6a1ReduceMod() int {
	if c.Base.Cons < 6 {
		return 1
	}
	return 0
}
