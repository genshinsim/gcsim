package kuki

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var chargeFrames []int
var chargeHitmarks = []int{14, 25}
var chargeHitlagHaltFrame = []float64{0, 0.10}
var chargeDefHalt = []bool{false, true}

func init() {
	chargeFrames = frames.InitAbilSlice(35) // CA -> N1/E/Q
	chargeFrames[action.ActionDash] = 31    // CA -> D
	chargeFrames[action.ActionJump] = 31    // CA -> J
	chargeFrames[action.ActionSwap] = 29    // CA -> Swap
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if c.Core.Player.LastAction.Type != action.ActionAttack {
		return action.Info{}, player.ErrInvalidChargeAction
	}

	for i, mult := range charge {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Charge %v", i),
			AttackTag:          attacks.AttackTagExtra,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeSlash,
			Element:            attributes.Physical,
			Durability:         25,
			Mult:               mult[c.TalentLvlAttack()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   chargeHitlagHaltFrame[i] * 60,
			CanBeDefenseHalted: chargeDefHalt[i],
		}
		// only the last multihit has hitlag so no need for char queue here
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 2.8),
			chargeHitmarks[i],
			chargeHitmarks[i],
		)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmarks[len(chargeHitmarks)-1],
		State:           action.ChargeAttackState,
	}, nil
}
