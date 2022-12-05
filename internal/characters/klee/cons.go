package klee

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c1(delay int) {
	if c.Base.Cons < 1 {
		return
	}
	//0.1 base change, + 0.08 every failure
	if c.Core.Rand.Float64() > c.c1Chance {
		//failed
		c.c1Chance += 0.08
		return
	}
	c.c1Chance = 0.1

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Sparks'n'Splash C1",
		AttackTag:          combat.AttackTagElementalBurst,
		ICDTag:             combat.ICDTagElementalBurst,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeDefault,
		Element:            attributes.Pyro,
		Durability:         25,
		Mult:               1.2 * burst[c.TalentLvlBurst()],
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.PrimaryTarget(), 1.5), 0, delay)
}

func (c *char) c2(a combat.AttackCB) {
	if c.Base.Cons < 2 {
		return
	}
	e, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	e.AddDefMod(enemy.DefMod{
		Base:  modifier.NewBaseWithHitlag("kleec2", 10*60),
		Value: -0.233,
	})
}
