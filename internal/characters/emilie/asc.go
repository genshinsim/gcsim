package emilie

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

const (
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
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.lumidoucePos, c.Core.Combat.PrimaryTarget(), nil, 3),
		a1Hitmark,
		a1Hitmark,
	)
}
