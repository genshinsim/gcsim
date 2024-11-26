package mavuika

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}
	c.nightsoulState.MaxPoints = 120
	c.fightingSpiritMult = 1.25
}

func (c *char) c1Atk() {
	if c.Base.Cons < 1 {
		return
	}
	buff := make([]float64, attributes.EndStatType)
	buff[attributes.ATKP] = 0.4
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag("mavuika-c1-atk", 8*60),
		Amount: func(a *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			return buff, true
		},
	})
}

func (c *char) c2BaseIncrease(activate bool) {
	if c.Base.Cons < 2 {
		return
	}

	if activate {
		c.Base.Atk += 300
	} else {
		c.Base.Atk -= 300
	}
}

func (c *char) c2FlatIncrease(tag attacks.AttackTag) float64 {
	if c.Base.Cons < 2 {
		return 0
	}
	switch tag {
	case attacks.AttackTagNormal:
		return 0.8 * c.TotalAtk()
	case attacks.AttackTagExtra:
		return 1.3 * c.TotalAtk()
	case attacks.AttackTagElementalBurst:
		return 1.8 * c.TotalAtk()
	default:
		return 0
	}
}

func (c *char) c2AddDefMod() {
	if c.Base.Cons < 2 {
		return
	}
	// TODO: the actual detection range
	enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 7), nil)
	for _, t := range enemies {
		// if you set duration=-1, it won't work
		t.AddDefMod(combat.DefMod{
			Base:  modifier.NewBaseWithHitlag("mavuika-c2-def-shred", 6000),
			Value: -0.2,
		})
	}
}

func (c *char) c2DeleteDefMod() {
	if c.Base.Cons < 2 {
		return
	}
	// all enemies
	enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 100), nil)
	for _, t := range enemies {
		if t.DefModIsActive("mavuika-c2-def-shred") {
			t.DeleteDefMod("mavuika-c2-def-shred")
		}
	}
}
