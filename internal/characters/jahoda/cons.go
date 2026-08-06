package jahoda

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c4Key = "jahoda-c4-flat-energy"
	c6Key = "jahoda-c6"
)

func (c *char) makeA1CB(a info.AttackCB) {
	if c.Base.Cons < 1 {
		return
	}

	// 50% to bounce
	if c.Core.Rand.Float64() < 0.5 {
		triggered := false

		if triggered {
			return
		}
		if a.Target.Type() != info.TargettableEnemy {
			return
		}
		triggered = true

		// default to bounce onto the original enemy
		var target info.Target = a.Target

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

		c.Core.QueueAttack(
			aiC1,
			combat.NewCircleHitOnTarget(target, nil, 4),
			0,
			c.meowballTravel+c1BounceHitmark,
			nil,
		)
	}

}

func (c *char) c4() {
	if c.Base.Cons < 4 {
		return
	}

	c.AddEnergy(c4Key, 4)
}

func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}

	c.c6Buff = make([]float64, attributes.EndStatType)
	c.c6Buff = make([]float64, attributes.EndStatType)

	c.c6Buff[attributes.CR] = 0.05
	c.c6Buff[attributes.CD] = 0.40

	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase(c6Key, 20*60),
			AffectedStat: attributes.CR,
			Amount: func() []float64 {
				if char.Moonsign < 1 {
					return nil
				}
				return c.c6Buff
			},
		})

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase(c6Key, 20*60),
			AffectedStat: attributes.CD,
			Amount: func() []float64 {
				if char.Moonsign < 1 {
					return nil
				}
				return c.c6Buff
			},
		})

		c.Core.Log.NewEvent("jahoda c6 triggered", glog.LogCharacterEvent, c.Index()).
			Write("cr", c.c6Buff[attributes.CR]).
			Write("cd", c.c6Buff[attributes.CD]).
			Write("expiry", c.Core.F+20*60)

	}
}
