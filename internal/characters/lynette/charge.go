package lynette

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var (
	chargeFrames []int
	// TODO: proper frames, currently using kirara
	chargeHitmarks = []int{20, 27}
)

func init() {
	// TODO: proper frames, currently using kirara
	chargeFrames = frames.InitAbilSlice(52) // CA -> Walk
	chargeFrames[action.ActionAttack] = 43
	chargeFrames[action.ActionSkill] = 43
	chargeFrames[action.ActionBurst] = 43
	chargeFrames[action.ActionDash] = 38
	chargeFrames[action.ActionJump] = 38
	chargeFrames[action.ActionSwap] = 37
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	for i, mult := range charge {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       fmt.Sprintf("Charge %v", i),
			AttackTag:  attacks.AttackTagExtra,
			ICDTag:     attacks.ICDTagNormalAttack,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeSlash,
			Element:    attributes.Physical,
			Durability: 25,
			Mult:       mult[c.TalentLvlAttack()],
		}
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: 0.8},
				2.2,
			),
			chargeHitmarks[i],
			chargeHitmarks[i],
		)
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmarks[len(chargeHitmarks)-1],
		State:           action.ChargeAttackState,
	}
}
