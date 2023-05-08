package qiqi

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
	chargeFrames   []int
	chargeHitmarks = []int{15, 29}
	chargeRadius   = []float64{2, 2.8}
	chargeOffsets  = []float64{0.5, 1}
)

func init() {
	chargeFrames = frames.InitAbilSlice(76) // CA -> N1
	chargeFrames[action.ActionSkill] = 67   // CA -> E
	chargeFrames[action.ActionBurst] = 67   // CA -> Q
	chargeFrames[action.ActionDash] = 32    // CA -> D
	chargeFrames[action.ActionJump] = 30    // CA -> J
	chargeFrames[action.ActionSwap] = 29    // CA -> Swap
}

// Standard charge attack
func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
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
			combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: chargeOffsets[i]},
				chargeRadius[i],
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
