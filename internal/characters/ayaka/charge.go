package ayaka

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int
var chargeHitmarks = []int{27, 33, 39}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		Abil:       "Charge",
		ActorIndex: c.Index,
		AttackTag:  combat.AttackTagExtra,
		ICDTag:     combat.ICDTagExtraAttack,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       ca[c.TalentLvlAttack()],
	}

	for i := 0; i < 3; i++ {
		c.Core.QueueAttack(
			ai,
			combat.NewDefCircHit(2, false, combat.TargettableEnemy),
			chargeHitmarks[i],
			chargeHitmarks[i],
			c.c1,
			c.c6,
		)
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmarks[len(chargeHitmarks)-1],
		Post:            chargeHitmarks[len(chargeHitmarks)-1],
		State:           action.ChargeAttackState,
	}
}
