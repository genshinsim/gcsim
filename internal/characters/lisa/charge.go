package lisa

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int

const (
	// hitmark frame, includes CA windup
	chargeHitmark = 70
	// TODO: stacks technically only last 15s and each stack has its own timer
	conductiveTag = "lisa-conductive-stacks"
)

func init() {
	chargeFrames = frames.InitAbilSlice(93)
	chargeFrames[action.ActionAttack] = 91
	chargeFrames[action.ActionCharge] = 90
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionSwap] = 90
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
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

	// skip CA windup if we're in NA or Skill (Hold) animation
	windup := 0
	switch c.Core.Player.CurrentState() {
	case action.NormalAttackState:
		windup = 14
	case action.SkillState:
		if c.Core.Player.LastAction.Param["hold"] != 0 {
			windup = 14
		}
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTargetFanAngle(
			c.Core.Combat.Player(),
			combat.Point{Y: 1},
			10,
			40,
		),
		chargeHitmark-windup,
		chargeHitmark-windup,
		c.makeA1CB(),
	)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return chargeFrames[next] - windup },
		AnimationLength: chargeFrames[action.InvalidAction] - windup,
		CanQueueAfter:   chargeHitmark - windup,
		State:           action.ChargeAttackState,
	}
}
