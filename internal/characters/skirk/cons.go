package skirk

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var c4Atkp = []float64{0.0, 0.1, 0.2, 0.4}

const c2Key = "skirk-c2"
const c6Dur = 15 * 60

func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Far to Fall",
		Mult:               5,
		AttackTag:          attacks.AttackTagExtra,
		ICDTag:             attacks.ICDTagExtraAttack,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeSlash,
		Element:            attributes.Cryo,
		Durability:         25,
		CanBeDefenseHalted: false,
	}
	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		nil,
		5,
	)
	c.Core.QueueAttack(ai, ap, 3, 3)
}

func (c *char) c2OnSkill() {
	if c.Base.Cons < 2 {
		return
	}
	c.AddSerpentsSubtlety(c.Base.Key.String()+"-c2", 10.0)
}

func (c *char) c2OnBurstRuin() float64 {
	if c.Base.Cons < 2 {
		return 0
	}
	return 10
}

func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}
	c.c2Atk = make([]float64, attributes.EndStatType)
	c.c2Atk[attributes.ATKP] = 0.7
}

func (c *char) c2OnBurstExtinction() {
	if c.Base.Cons < 2 {
		return
	}

	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(c2Key, 12.5*60),
		AffectedStat: attributes.ATKP,
		Amount: func() ([]float64, bool) {
			if !c.StatusIsActive(skillKey) {
				return nil, false
			}
			return c.c2Atk, true
		},
	})
}

func (c *char) c4Init() {
	if c.Base.Cons < 4 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("skirk-c4", -1),
		AffectedStat: attributes.ATKP,
		Amount: func() ([]float64, bool) {
			m[attributes.ATKP] = c4Atkp[c.getA4Stacks()]
			return m, true
		},
	})
}

func (c *char) c6OnVoidAbsorb() {
	if c.Base.Cons < 6 {
		return
	}
	c.c6Stacks.PushOverwrite(c.TimePassed)
}

func (c *char) c6OnBurstRuin() {
	if c.Base.Cons < 6 {
		return
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Havoc: Sever (Burst)",
		Mult:       7.5,
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Cryo,
		Durability: 25,
	}
	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		nil,
		5,
	)

	for !c.c6Stacks.IsEmpty() {
		src, _ := c.c6Stacks.Pop()
		if src+c6Dur < c.TimePassed {
			continue
		}
		c.Core.QueueAttack(ai, ap, 3, 3)
	}
}

func (c *char) c6OnAttackCB() func(a combat.AttackCB) {
	if c.Base.Cons < 6 {
		return nil
	}

	switch c.NormalCounter {
	case 2, 4:
	default:
		return nil
	}

	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}

		hasStack := false
		for !c.c6Stacks.IsEmpty() {
			src, _ := c.c6Stacks.Pop()
			if src+c6Dur >= c.TimePassed {
				hasStack = true
				break
			}
		}
		if !hasStack {
			return
		}

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Havoc: Sever (Normal)",
			Mult:       1.8,
			AttackTag:  attacks.AttackTagNormal,
			ICDTag:     attacks.ICDTagNormalAttack,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeSlash,
			Element:    attributes.Cryo,
			Durability: 25,
		}
		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			nil,
			5,
		)
		for i := 1; i <= 3; i++ {
			c.Core.QueueAttack(ai, ap, i*3, i*3)
		}
	}
}

func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}
	c.c6Stacks = NewRingQueue[int](3)

	// TODO: OnCharacterHurt?
}
