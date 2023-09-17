package kazuha

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var chargeFrames []int
var chargeHitmarks = []int{20, 20}

func init() {
	chargeFrames = frames.InitAbilSlice(55)
	chargeFrames[action.ActionSkill] = 30
	chargeFrames[action.ActionBurst] = 30
	chargeFrames[action.ActionDash] = 20
	chargeFrames[action.ActionJump] = 20
	chargeFrames[action.ActionSwap] = 28
}

func (c *char) ChargeAttack(p map[string]int) action.Info {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNormalAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Physical,
		Durability: 25,
	}

	for i, mult := range charge {
		ai.Mult = mult[c.TalentLvlAttack()]
		ai.Abil = fmt.Sprintf("Charge %v", i)
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, 2.2),
			chargeHitmarks[i],
			chargeHitmarks[i],
		)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmarks[len(chargeHitmarks)-1],
		State:           action.ChargeAttackState,
	}
}
