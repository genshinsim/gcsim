package ayaka

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int
var chargeHitmarks = []int{27, 33, 39}

func init() {
	chargeFrames = frames.InitAbilSlice(71)
	chargeFrames[action.ActionSkill] = 62
	chargeFrames[action.ActionBurst] = 63
	chargeFrames[action.ActionDash] = chargeHitmarks[len(chargeHitmarks)-1]
	chargeFrames[action.ActionJump] = chargeHitmarks[len(chargeHitmarks)-1]
	chargeFrames[action.ActionSwap] = chargeHitmarks[len(chargeHitmarks)-1]
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		Abil:       "Charge",
		ActorIndex: c.Index,
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagExtraAttack,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       ca[c.TalentLvlAttack()],
	}

	for i := 0; i < 3; i++ {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				nil,
				1,
			),
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
		State:           action.ChargeAttackState,
	}
}
