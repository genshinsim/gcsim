package skirk

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var c4Atkp = []float64{0.0, 0.1, 0.2, 0.4}

const (
	c2Key = "skirk-c2"
	c6Dur = 15 * 60
)

func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}
	ai := info.AttackInfo{
		ActorIndex:         c.Index(),
		Abil:               "Far to Fall",
		Mult:               5,
		AttackTag:          attacks.AttackTagExtra,
		ICDTag:             attacks.ICDTagSkirkCons,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeDefault,
		Element:            attributes.Cryo,
		Durability:         25,
		CanBeDefenseHalted: false,
	}
	ap := combat.NewBoxHitOnTarget(
		c.Core.Combat.Player(),
		nil,
		4,
		4,
	)
	// TODO: Attack Delay on C1?
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
	// delay buff to the end of burst. c1 hits from burst don't benefit from c2
	c.QueueCharTask(func() {
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
	}, 39)
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
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Havoc: Sever (Burst)",
		Mult:       7.5 * c.a4MultBurst(),
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
	}
	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		nil,
		5,
	)

	filter := func(src int) bool {
		return src+c6Dur >= c.TimePassed
	}

	count := c.c6Stacks.Count(filter)
	c.c6Stacks.Clear()
	for range count {
		// TODO: Record Attack Delay on C6?
		c.Core.QueueAttack(ai, ap, 0.2*60, 0.2*60)
	}
}

func (c *char) c6OnAttackCB() func(a info.AttackCB) {
	if c.Base.Cons < 6 {
		return nil
	}

	switch c.NormalCounter {
	case 2, 4:
	default:
		return nil
	}

	return func(a info.AttackCB) {
		if a.Target.Type() != info.TargettableEnemy {
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

		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Havoc: Sever (Normal)",
			Mult:       1.8 * c.a4MultAttack(),
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
