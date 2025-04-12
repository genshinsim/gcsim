package mizuki

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int

// hitmark frame, includes CA windup
const chargeHitmark = 39

func init() {
	chargeFrames = frames.InitAbilSlice(62)
	chargeFrames[action.ActionDash] = 37
	chargeFrames[action.ActionJump] = 37
	chargeFrames[action.ActionWalk] = 81
	chargeFrames[action.ActionSwap] = 64
}

// Standard charge attack
// CA has no travel time
func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	// CA windup on idle
	windup := 14
	if c.Core.Player.CurrentState() == action.Idle {
		windup = 0
	}

	radius := 3.5

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, radius),
		chargeHitmark-windup,
		chargeHitmark-windup,
	)

	return action.Info{
		Frames:          func(next action.Action) int { return chargeFrames[next] - windup },
		AnimationLength: chargeFrames[action.InvalidAction] - windup,
		CanQueueAfter:   chargeHitmark - windup,
		State:           action.ChargeAttackState,
	}, nil
}
