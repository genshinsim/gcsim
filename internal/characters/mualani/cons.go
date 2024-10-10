package mualani

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c4key = "mualani-c4"

func (c *char) c1() float64 {
	if c.Base.Cons < 1 {
		return 0.0
	}
	if c.c1Done {
		return 0.0
	}
	c.c1Done = true
	c.c6()
	return 0.66
}

func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}
	c.momentumStacks = 2
}

func (c *char) c2puffer() {
	if c.Base.Cons < 2 {
		return
	}
	if c.Base.Ascension < 1 {
		return
	}

	c.momentumStacks = min(c.momentumStacks+1, 3)
	if c.a1Count == 2 {
		for i := 0; i < 12; i++ {
			c.QueueCharTask(func() {
				c.nightsoulState.GeneratePoints(1)
			}, i*10+10)
		}
	}
}
func (c *char) c4() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.75
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag(c4key, -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag == attacks.AttackTagElementalBurst {
				return m, true
			}
			return nil, false
		},
	})
}

func (c *char) c4puffer() {
	if c.Base.Cons < 4 {
		return
	}
	if c.Base.Ascension < 1 {
		return
	}

	c.AddEnergy(c4key, 8)
}

func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}
	c.c1Done = false
}
