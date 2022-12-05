package albedo

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int
var chargeHitmarks = []int{20, 20} // CA-1 and CA-2 hit at the same time

func init() {
	chargeFrames = frames.InitAbilSlice(56)                                 // CA -> N1
	chargeFrames[action.ActionSkill] = 34                                   // CA -> E
	chargeFrames[action.ActionBurst] = 34                                   // CA -> Q
	chargeFrames[action.ActionDash] = chargeHitmarks[len(chargeHitmarks)-1] // CA -> D
	chargeFrames[action.ActionJump] = chargeHitmarks[len(chargeHitmarks)-1] // CA -> J
	chargeFrames[action.ActionSwap] = 33                                    // CA -> Swap
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		AttackTag:  combat.AttackTagNormal,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeSlash,
		Element:    attributes.Physical,
		Durability: 25,
	}

	for i, mult := range charge {
		ai.Mult = mult[c.TalentLvlAttack()]
		ai.Abil = fmt.Sprintf("Charge %v", i)
		c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 2.2), chargeHitmarks[i], chargeHitmarks[i])
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmarks[len(chargeHitmarks)-1],
		State:           action.ChargeAttackState,
	}
}
