package xiangling

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int

const chargeHitmark = 24

func init() {
	chargeFrames = frames.InitAbilSlice(69)
	chargeFrames[action.ActionAttack] = 67
	chargeFrames[action.ActionBurst] = 67
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionSwap] = 66
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Charge",
		AttackTag:          combat.AttackTagExtra,
		ICDTag:             combat.ICDTagExtraAttack,
		ICDGroup:           combat.ICDGroupPoleExtraAttack,
		Element:            attributes.Physical,
		Durability:         25,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
		IsDeployable:       true,
		Mult:               nc[c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 0.1), chargeHitmark, chargeHitmark)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmark,
		State:           action.ChargeAttackState,
	}
}
