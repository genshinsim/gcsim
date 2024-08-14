package emilie

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a4ModKey = "emilie-a4"

	a1Hitmark = 16
)

func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	c.SetTag(lumidouceScent, 0)

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Cleardew Cologne (A1)",
		AttackTag:  attacks.AttackTagNone,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       6,
	}
	c.applyC6Bonus(&ai)
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.lumidoucePos, c.Core.Combat.PrimaryTarget(), nil, 3),
		a1Hitmark,
		a1Hitmark,
		c.c2,
	)
}

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(a4ModKey, -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			x, ok := t.(*enemy.Enemy)
			if !ok {
				return nil, false
			}
			if !x.IsBurning() {
				return nil, false
			}
			m[attributes.DmgP] = c.TotalAtk() / 1000 * 0.15
			return m, true
		},
	})
}
