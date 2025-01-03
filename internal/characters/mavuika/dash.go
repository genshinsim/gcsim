package mavuika

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

func (c *char) Dash(p map[string]int) (action.Info, error) {
	if c.flamestriderModeActive && c.nightsoulState.HasBlessing() {
		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Flamestrider Sprint",
			AttackTag:      attacks.AttackTagNone,
			ICDTag:         attacks.ICDTagExtraAttack,
			AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypeBlunt,
			PoiseDMG:       75.0,
			Element:        attributes.Pyro,
			Durability:     25,
			Mult:           1.485,
		}
		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: 1.0},
			1.2,
		)
		c.Core.QueueAttack(ai, ap, 10, 10)
		c.reduceNightsoulPoints(10)
	}

	return c.Character.Dash(p)
}
