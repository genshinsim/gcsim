package kirara

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

// based on keqing frames
// TODO: update frames
var (
	chargeFrames   []int
	chargeHitmarks = []int{12, 34, 36}
	chargeRadius   = []float64{1.5, 2.2, 2.3}
	chargeOffsets  = []float64{1.5, 1.5, 1.8}
)

func init() {
	chargeFrames = frames.InitAbilSlice(48)
	chargeFrames[action.ActionSkill] = 47
	chargeFrames[action.ActionBurst] = 47
	chargeFrames[action.ActionDash] = chargeHitmarks[len(chargeHitmarks)-1]
	chargeFrames[action.ActionJump] = chargeHitmarks[len(chargeHitmarks)-1]
	chargeFrames[action.ActionSwap] = chargeHitmarks[len(chargeHitmarks)-1]
}

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
