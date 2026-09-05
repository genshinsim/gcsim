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

func (c *char) c4OnBurst(buff float64) {
	if c.Base.Cons < 4 {
		return
	}
	buff *= 0.5
	for _, char := range c.Core.Player.Chars() {
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

func (c *char) c4Init() {
	if c.Base.Cons < 4 {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		e, ok := args[0].(*enemy.Enemy)
		if !ok {
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
