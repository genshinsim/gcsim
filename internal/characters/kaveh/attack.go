package kaveh

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
	attackFrames          [][]int
	attackHitmarks        = []int{27, 22, 33, 40}
	attackHitlagHaltFrame = []float64{.1, .09, .09, .1}
	attackHitboxes        = [][]float64{{2.2, 2.2, 2.3, 2.3}, {3.1, 3.1, 3.2, 3.1}}
	attackOffsets         = []float64{0.5, 0.5, 0.5, 2.5}
	attackFanAngles       = []float64{260, 260, 240, 360}
)

func init() {
	attackFrames = make([][]int, len(attackHitmarks))
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 44)
	attackFrames[0][action.ActionAttack] = 35

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 45)
	attackFrames[1][action.ActionAttack] = 35

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 56)
	attackFrames[2][action.ActionAttack] = 42

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 81)
	attackFrames[3][action.ActionAttack] = 67
}

func (c *char) Attack(p map[string]int) action.Info {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:          attacks.AttackTagNormal,
		ICDTag:             attacks.ICDTagNormalAttack,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		Element:            attributes.Physical,
		Durability:         25,
		Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter] * 60,
		CanBeDefenseHalted: true,
	}

	// check burst status for hitbox
	attackIndex := 0
	if c.StatModIsActive(burstKey) {
		attackIndex = 1
	}
	ap := combat.NewCircleHitOnTargetFanAngle(
		c.Core.Combat.Player(),
		geometry.Point{Y: attackOffsets[c.NormalCounter]},
		attackHitboxes[attackIndex][c.NormalCounter],
		attackFanAngles[c.NormalCounter],
	)

	c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter], attackHitmarks[c.NormalCounter])

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}
