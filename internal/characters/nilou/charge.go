package nilou

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int

const chargeHitmark = 14

// TODO: cancel frames & hitlags
func init() {
	chargeFrames = frames.InitAbilSlice(14)
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	for i, mult := range charge {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       fmt.Sprintf("Charge %v", i),
			AttackTag:  combat.AttackTagExtra,
			ICDTag:     combat.ICDTagExtraAttack,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeSlash,
			Element:    attributes.Physical,
			Durability: 25,
			Mult:       mult[c.TalentLvlAttack()],
		}
		// only the last multihit has hitlag so no need for char queue here
		c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 0.5, false, combat.TargettableEnemy), chargeHitmark, chargeHitmark)
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmark,
		State:           action.ChargeAttackState,
	}
}
