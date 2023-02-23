package yaemiko

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int

const chargeHitmark = 90

func init() {
	chargeFrames = frames.InitAbilSlice(96) // CA -> N1/E/Q
	chargeFrames[action.ActionCharge] = 95  // CA -> CA
	chargeFrames[action.ActionDash] = 46    // CA -> D
	chargeFrames[action.ActionJump] = 47    // CA -> J
	chargeFrames[action.ActionSwap] = 94    // CA -> Swap
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	// Supposed to be ICDTagYaeCharged and ICDGroupYaeCharged. However, it's
	// essentially no ICD because it takes ~36f for the charge to hit again.
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	// skip CA windup if we're in NA animation
	windup := 0
	if c.Core.Player.CurrentState() == action.NormalAttackState {
		windup = 14
	}

	// TODO: check snapshot delay
	c.Core.QueueAttack(
		ai,
		combat.NewBoxHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 2, 2),
		0,
		chargeHitmark-windup,
	)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return chargeFrames[next] - windup },
		AnimationLength: chargeFrames[action.InvalidAction] - windup,
		CanQueueAfter:   chargeFrames[action.ActionDash] - windup, // earliest cancel is before hitmark
		State:           action.ChargeAttackState,
	}
}
