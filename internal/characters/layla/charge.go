package layla

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int
var chargeHitmarks = []int{16, 16 + 11}

func init() {
	chargeFrames = frames.InitAbilSlice(49) // CA -> N1/W
	chargeFrames[action.ActionSkill] = 34   // CA -> E
	chargeFrames[action.ActionBurst] = 34   // CA -> Q
	chargeFrames[action.ActionDash] = 27    // CA -> D
	chargeFrames[action.ActionJump] = 27    // CA -> J
	chargeFrames[action.ActionSwap] = 29    // CA -> Swap
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		AttackTag:  combat.AttackTagExtra,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeSlash,
		Element:    attributes.Physical,
		Durability: 25,
	}

	for i, mult := range charge {
		ai.Mult = mult[c.TalentLvlAttack()]
		ai.Abil = fmt.Sprintf("Charge %v", i)

		if i == 1 {
			ai.HitlagFactor = 0.01
			ai.HitlagHaltFrames = 0.06 * 60
			ai.CanBeDefenseHalted = true
		}

		c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 2.8), chargeHitmarks[i], chargeHitmarks[i])
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeFrames[action.ActionJump], // earliest cancel
		State:           action.ChargeAttackState,
	}
}
