package yunjin

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	attackFrames          [][]int
	attackHitmarks        = [][]int{{15}, {13}, {8, 23}, {11, 23}, {15}}
	attackHitlagHaltFrame = [][]float64{{0.03}, {0.03}, {0, 0.03}, {0, 0.03}, {0.04}}
	attackDefHalt         = [][]bool{{true}, {true}, {false, true}, {false, false}, {true}}
	attackHitboxes        = [][]float64{{2}, {2}, {2}, {2, 3.2}, {2.4}}
	attackOffsets         = []float64{0.8, 0.8, 0.8, -1.2, 1.1}
	attackFanAngles       = []float64{220, 220, 220, 360, 360}
)

const normalHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 25) // N1 -> CA
	attackFrames[0][action.ActionAttack] = 20                                // N1 -> N2

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 29) // N2 -> CA
	attackFrames[1][action.ActionAttack] = 22                                // N1 -> N2

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 32) // N3 -> CA
	attackFrames[2][action.ActionAttack] = 31                                // N1 -> N2

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][1], 46) // N4 -> CA
	attackFrames[3][action.ActionAttack] = 45                                // N1 -> N2

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 67) // N5 -> N1
	attackFrames[4][action.ActionCharge] = 500                               // N5 -> CA, TODO: this action is illegal; need better way to handle it
}

// Normal attack damage queue generator
// Very standard
func (c *char) Attack(p map[string]int) action.ActionInfo {
	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             combat.ICDTagNormalAttack,
			ICDGroup:           combat.ICDGroupDefault,
			StrikeType:         combat.StrikeTypeSlash,
			Element:            attributes.Physical,
			Durability:         25,
			Mult:               mult[c.TalentLvlAttack()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}
		ap := combat.NewCircleHitOnTargetFanAngle(
			c.Core.Combat.Player(),
			combat.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
			attackFanAngles[c.NormalCounter],
		)
		if c.NormalCounter == 3 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				combat.Point{Y: attackOffsets[c.NormalCounter]},
				attackHitboxes[c.NormalCounter][0],
				attackHitboxes[c.NormalCounter][1],
			)
		}
		c.QueueCharTask(func() {
			c.Core.QueueAttack(ai, ap, 0, 0)
		}, attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}
