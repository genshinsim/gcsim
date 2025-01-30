package mavuika

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c6IcdKey = "mavuika-c6-icd"
const c1Key = "mavuika-c1"

func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		c.nightsoulState.MaxPoints = 80
		return
	}
	c.nightsoulState.MaxPoints = 120
	c.c1buff = make([]float64, attributes.EndStatType)
	c.c1buff[attributes.ATKP] = 0.4
}
func (c *char) c1FightingSpiritEff() float64 {
	if c.Base.Cons < 1 {
		return 1.0
	}
	return 1.25
}
func (c *char) c1OnFightingSpirit() {
	if c.Base.Cons < 1 {
		return
	}
	c.AddStatMod(character.StatMod{
		Base: modifier.NewBaseWithHitlag(c1Key, 10*60),
		Amount: func() ([]float64, bool) {
			return c.c1buff, true
		},
	})
}

func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.BaseATK] = 200
	c.AddStatMod(character.StatMod{
		Base: modifier.NewBase("mavuika-c2-base-atk", -1),
		Amount: func() ([]float64, bool) {
			if c.nightsoulState.HasBlessing() {
				return m, true
			}
			return nil, false
		},
	})
}

func (c *char) c2Ring() {
	if c.Base.Cons < 2 {
		return
	}
	if !c.isRingFollowing() {
		return
	}
	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		geometry.Point{Y: 1.0},
		6,
	)
	for _, e := range c.Core.Combat.EnemiesWithinArea(ap, nil) {
		e.AddDefMod(combat.DefMod{
			Base:  modifier.NewBaseWithHitlag("mavuika-c2", 30),
			Value: -0.2,
		})
	}
	c.QueueCharTask(c.c2Ring, 18)
}

func (c *char) c2BikeNA() float64 {
	if c.Base.Cons < 2 {
		return 0.0
	}
	if c.armamentState != bike {
		return 0.0
	}
	return 0.6 * c.TotalAtk()
}
func (c *char) c2BikeCA() float64 {
	if c.Base.Cons < 2 {
		return 0.0
	}
	if c.armamentState != bike {
		return 0.0
	}
	return 0.9 * c.TotalAtk()
}

func (c *char) c2BikeQ() float64 {
	if c.Base.Cons < 2 {
		return 0.0
	}
	if c.armamentState != bike {
		return 0.0
	}
	return 1.2 * c.TotalAtk()
}

func (c *char) c4BonusVal() float64 {
	if c.Base.Cons < 4 {
		return 0.0
	}
	return 0.1
}

func (c *char) c4DecayRate() float64 {
	if c.Base.Cons < 4 {
		return 1.0 / (20 * 60)
	}
	return 0.0
}

// this is just used for c2
func (c *char) isRingFollowing() bool {
	if c.Base.Cons < 6 {
		return c.armamentState == ring
	}
	return true
}

func (c *char) c6RingCB() func(a combat.AttackCB) {
	if c.Base.Cons < 6 {
		return nil
	}
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if c.StatusIsActive(c6IcdKey) {
			return
		}
		c.AddStatus(c6IcdKey, 0.5*60, true)
		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Flamestrider (C6)",
			AttackTag:      attacks.AttackTagElementalArt,
			ICDTag:         attacks.ICDTagNone,
			AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypeBlunt,
			PoiseDMG:       75,
			Element:        attributes.Pyro,
			Durability:     0,
			Mult:           2.0,
		}
		ap := combat.NewCircleHitOnTarget(
			a.Target,
			nil,
			6,
		)
		c.Core.QueueAttack(ai, ap, 3, 3)
	}
}

func (c *char) c6Bike() {
	if c.Base.Cons < 6 {
		return
	}
	c.c6Src = c.Core.F
	c.QueueCharTask(c.c6RingAtk(c.c6Src), 180)
}

func (c *char) c6RingAtk(src int) func() {
	return func() {
		if c.c6Src != src {
			return
		}
		if c.armamentState != bike {
			return
		}
		if !c.nightsoulState.HasBlessing() {
			return
		}
		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Rings of Searing Radiance (C6)",
			AttackTag:      attacks.AttackTagElementalArt,
			ICDTag:         attacks.ICDTagNone,
			AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypePierce,
			Element:        attributes.Pyro,
			Durability:     0,
			Mult:           4.0,
		}
		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: 1.0},
			6,
		)
		c.Core.QueueAttack(ai, ap, 0, 0)
		c.QueueCharTask(c.c6RingAtk(src), 180)
	}
}
