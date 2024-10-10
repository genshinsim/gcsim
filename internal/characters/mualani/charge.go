package mualani

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var chargeFrames []int

const shortChargeHitmark = 71

func init() {
	chargeFrames = frames.InitAbilSlice(100)
	chargeFrames[action.ActionAttack] = 73
	chargeFrames[action.ActionCharge] = 85
	chargeFrames[action.ActionSkill] = 72
	chargeFrames[action.ActionBurst] = 73
	chargeFrames[action.ActionDash] = 72
	chargeFrames[action.ActionJump] = 73
	chargeFrames[action.ActionWalk] = 100
	chargeFrames[action.ActionSwap] = 71
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	// there is a windup out of dash/jump/walk/swap. Otherwise it is rolled into the Q/E/CA/NA -> CA frames
	windup := 0
	switch c.Core.Player.CurrentState() {
	case action.Idle, action.DashState, action.JumpState, action.WalkState, action.SwapState:
		windup = 13
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, 3.5)
	// TODO: Not sure of snapshot timing
	c.Core.QueueAttack(
		ai,
		ap,
		shortChargeHitmark+windup,
		shortChargeHitmark+windup,
	)

	return action.Info{
		Frames:          func(next action.Action) int { return windup + chargeFrames[next] },
		AnimationLength: windup + chargeFrames[action.InvalidAction],
		CanQueueAfter:   windup + chargeFrames[action.ActionDash],
		State:           action.ChargeAttackState,
	}, nil
}
