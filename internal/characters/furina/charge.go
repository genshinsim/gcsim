package furina

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var (
	chargeFrames  []int
	chargeHitmark = 30
	chargeOffset  = 0.0
)

func init() {
	chargeFrames = frames.InitAbilSlice(52) // C -> Walk
	chargeFrames[action.ActionAttack] = 43
	chargeFrames[action.ActionSkill] = 43
	chargeFrames[action.ActionBurst] = 43
	chargeFrames[action.ActionDash] = 38
	chargeFrames[action.ActionJump] = 38
	chargeFrames[action.ActionSwap] = 37
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Charge"),
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNormalAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: chargeOffset},
			2.8,
		),
		chargeHitmark,
		chargeHitmark,
	)
	if c.arkhe == ousia {
		c.arkhe = pneuma
	} else {
		c.arkhe = ousia
	}
	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmark,
		State:           action.ChargeAttackState,
	}, nil
}
