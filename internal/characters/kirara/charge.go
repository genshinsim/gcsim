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

var (
	chargeFrames   []int
	chargeHitmarks = []int{20, 27, 37}
	chargeOffsets  = []float64{1, 1, 1.5}
)

func init() {
	chargeFrames = frames.InitAbilSlice(52) // C -> Walk
	chargeFrames[action.ActionAttack] = 43
	chargeFrames[action.ActionSkill] = 43
	chargeFrames[action.ActionBurst] = 43
	chargeFrames[action.ActionDash] = 38
	chargeFrames[action.ActionJump] = 38
	chargeFrames[action.ActionSwap] = 37
}

func (c *char) ChargeAttack(p map[string]int) action.Info {
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
		if i == 2 {
			ai.HitlagFactor = 0.01
			ai.HitlagHaltFrames = 0.1 * 60
			ai.CanBeDefenseHalted = true
		}
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: chargeOffsets[i]},
				2.8,
			),
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
