package sucrose

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int

// hitmark frame, includes CA windup
const chargeHitmark = 54

func init() {
	chargeFrames = frames.InitAbilSlice(69)
	chargeFrames[action.ActionCharge] = 66
	chargeFrames[action.ActionSkill] = 60
	chargeFrames[action.ActionBurst] = 61
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionSwap] = chargeHitmark // idk if this is correct or not
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  combat.AttackTagExtra,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	// skip CA windup if we're in NA/CA animation
	windup := 0
	switch c.Core.Player.CurrentState() {
	case action.NormalAttackState, action.ChargeAttackState:
		windup = 15
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: -0.2}, 3.2, 7.5),
		chargeHitmark-windup,
		chargeHitmark-windup,
	)

	if c.Base.Cons >= 4 {
		c.c4()
	}

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return chargeFrames[next] - windup },
		AnimationLength: chargeFrames[action.InvalidAction] - windup,
		CanQueueAfter:   chargeHitmark - windup,
		State:           action.ChargeAttackState,
	}
}
