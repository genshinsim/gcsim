package nahida

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int

const chargeHitmark = 65

func init() {
	chargeFrames = frames.InitAbilSlice(65)
	chargeFrames[action.ActionDash] = 21
	chargeFrames[action.ActionJump] = 21
	chargeFrames[action.ActionWalk] = 69
	chargeFrames[action.ActionSwap] = 21

}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	windup := 0
	switch c.Core.Player.CurrentState() {
	case action.NormalAttackState:
		windup = 10
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			combat.Point{Y: -3},
			6,
			6,
		),
		chargeHitmark-windup,
		chargeHitmark-windup,
		c.makeC6CB(),
	)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return chargeFrames[next] - windup },
		AnimationLength: chargeFrames[action.InvalidAction] - windup,
		CanQueueAfter:   chargeFrames[action.ActionDash] - windup,
		State:           action.ChargeAttackState,
	}
}
