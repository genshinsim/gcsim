package chiori

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
	chargeFrames          []int
	chargeHitmarks        = []int{25, 26}
	chargeHitlagHaltFrame = []float64{0, 0.06}
	chargeDefHalt         = []bool{false, true}
	chargeRadius          = []float64{2.3, 2.4}
	chargeOffsets         = []float64{1.4, 1.6}
)

func init() {
	chargeFrames = frames.InitAbilSlice(44) // CA -> Walk
	chargeFrames[action.ActionAttack] = 39
	chargeFrames[action.ActionSkill] = 39
	chargeFrames[action.ActionBurst] = 39
	chargeFrames[action.ActionDash] = chargeHitmarks[len(chargeHitmarks)-1]
	chargeFrames[action.ActionJump] = chargeHitmarks[len(chargeHitmarks)-1]
	chargeFrames[action.ActionSwap] = chargeHitmarks[len(chargeHitmarks)-1]
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:   c.Index,
		AttackTag:    attacks.AttackTagExtra,
		ICDTag:       attacks.ICDTagNormalAttack,
		ICDGroup:     attacks.ICDGroupDefault,
		StrikeType:   attacks.StrikeTypeSlash,
		Element:      attributes.Physical,
		Durability:   25,
		HitlagFactor: 0.01,
	}

	for i, mult := range charge {
		ai.Mult = mult[c.TalentLvlAttack()]
		ai.Abil = fmt.Sprintf("Charge %v", i)
		ai.HitlagHaltFrames = chargeHitlagHaltFrame[i] * 60
		ai.CanBeDefenseHalted = chargeDefHalt[i]
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: chargeOffsets[i]}, chargeRadius[i]),
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
