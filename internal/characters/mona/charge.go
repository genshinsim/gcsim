package mona

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int

const chargeHitmark = 66

func init() {
	chargeFrames = frames.InitAbilSlice(113) // CA -> N1
	chargeFrames[action.ActionCharge] = 105  // CA -> CA
	chargeFrames[action.ActionSkill] = 92    // CA -> E
	chargeFrames[action.ActionBurst] = 68    // CA -> Q
	chargeFrames[action.ActionDash] = 72     // CA -> D
	chargeFrames[action.ActionJump] = 64     // CA -> J
	chargeFrames[action.ActionSwap] = 54     // CA -> Swap
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  combat.AttackTagExtra,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	// add windup if we're in idle or swap only
	// TODO: this ignores N4 -> CA (which should be illegal anyways)
	windup := 14
	if c.Core.Player.CurrentState() == action.Idle || c.Core.Player.CurrentState() == action.SwapState {
		windup = 0
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			nil,
			3,
		),
		chargeHitmark-windup,
		chargeHitmark-windup,
	)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return chargeFrames[next] - windup },
		AnimationLength: chargeFrames[action.InvalidAction] - windup,
		CanQueueAfter:   chargeFrames[action.ActionSwap] - windup, // earliest cancel is before hitmark
		State:           action.ChargeAttackState,
	}
}
