package albedo

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int
var chargeHitmarks = []int{39, 49}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		AttackTag:  combat.AttackTagNormal,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Physical,
		Durability: 25,
	}

	//TODO: damage frame
	for i, mult := range charge {
		ai.Mult = mult[c.TalentLvlAttack()]
		ai.Abil = fmt.Sprintf("Charge %v", i)
		c.Core.QueueAttack(ai, combat.NewDefCircHit(0.1, false, combat.TargettableEnemy), chargeHitmarks[i], chargeHitmarks[i])
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmarks[len(chargeHitmarks)-1],
		Post:            chargeHitmarks[len(chargeHitmarks)-1],
		State:           action.ChargeAttackState,
	}
}
