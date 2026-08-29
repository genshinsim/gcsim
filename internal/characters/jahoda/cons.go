package jahoda

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c4Key   = "jahoda-c4-flat-energy"
	c6Key = "jahoda-c6"
)

func (c *char) makeC1CB(a info.AttackCB) {
	if c.Base.Cons < 1 {
		return
	}

	if a.Target.Type() != info.TargettableEnemy {
		return
	}

	// 50% to bounce
	if c.Core.Rand.Float64() < 0.5 {
		return
	}

	// default to bounce onto the original enemy
	target := a.Target

	// prefer a different nearby enemy when one exists, the exact radius of detection is unknown
	next := c.Core.Combat.RandomEnemyWithinArea(combat.NewCircleHitOnTarget(a.Target, nil, 8), func(t info.Enemy) bool {
		return t.Key() != a.Target.Key()
	},
	)

	if next != nil {
		target = next
	}

	// queue the attack
	aiC1 := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Meowball (C1)",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagJahodaCons,
		ICDGroup:   attacks.ICDGroupJahodaCons,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    c.flaskAbsorb,
		Durability: 25,
		Mult:       meowball[c.TalentLvlSkill()],
	}

	// TODO: C1 use the same snapshot as original meowball, not tested
	c.Core.QueueAttack(
		aiC1,
		combat.NewCircleHitOnTarget(target, nil, 4),
		0,
		c.meowballTravel+c1BounceHitmark,
		nil,
	)
}

func (c *char) c2Init(eleCountMap map[attributes.Element]int) {
	if c.Base.Cons < 2 {
		return
	}

	if c.Core.Player.GetMoonsignLevel() < 2 {
		return
	}

	secondHighestEleCount := 0

	for _, ele := range elePriority {
		if ele == c.a1HighestEle {
			continue
		}

		if eleCountMap[ele] > secondHighestEleCount {
			secondHighestEleCount = eleCountMap[ele]
			c.c2NextHighestEle = ele
		}
	}

	if secondHighestEleCount == 0 {
		c.c2NextHighestEle = attributes.NoElement
	}
}

func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}

	if c.Core.Player.GetMoonsignLevel() < 2 {
		return
	}

	c.applyA1Buff(c.c2NextHighestEle)
}

func (c *char) c4() {
	if c.Base.Cons < 4 {
		return
	}

	c.AddEnergy(c4Key, 4)
}

func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}

	if c.Core.Player.GetMoonsignLevel() < 2 {
		return
	}

	c.c6Buff = make([]float64, attributes.EndStatType)
	c.c6Buff[attributes.CR] = 0.05
	c.c6Buff[attributes.CD] = 0.40
}

func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}

	if c.Core.Player.GetMoonsignLevel() < 2 {
		return
	}

	for _, char := range c.Core.Player.Chars() {
		if char.Moonsign < 1 {
			continue
		}

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(c6Key, 20*60),
			Amount: func() []float64 {
				return c.c6Buff
			},
		})
	}
}
