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
	attackFrames          [][]int
	attackHitmarks        = [][]int{{11}, {10}, {9, 20}, {8, 18}, {16}}
	attackHitlagHaltFrame = [][]float64{{0.03}, {0.03}, {0.03, 0.03}, {0.03, 0.03}, {0.12}}
	attackHitlagFactor    = [][]float64{{0.01}, {0.01}, {0.05, 0.05}, {0.05, 0.05}, {0.01}}
	attackHitboxes        = [][]float64{{1.2}, {1.4, 2.2}, {1.6}, {1.6}, {2.2}}
	attackOffsets         = [][]float64{{0.8}, {0}, {1, 0.6}, {0.6, 0.6}, {1}}
	attackFanAngles       = [][]float64{{360}, {360}, {30, 360}, {360, 360}, {360}}
)

const normalHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 21) // N1 -> CA
	attackFrames[0][action.ActionAttack] = 14                                // N1 -> N2

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 22) // N2 -> CA
	attackFrames[1][action.ActionAttack] = 21                                // N2 -> N3

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 33) // N3 -> N4/CA

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][1], 28) // N4 -> N5
	attackFrames[3][action.ActionCharge] = 26                                // N4 -> CA

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 53) // N5 -> N1
	attackFrames[4][action.ActionCharge] = 500                               // N5 -> CA, TODO: this action is illegal; need better way to handle it

}

// Standard attack - nothing special
func (c *char) Attack(p map[string]int) action.ActionInfo {

	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			Mult:               mult[c.TalentLvlAttack()],
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeSlash,
			Element:            attributes.Physical,
			Durability:         25,
			HitlagFactor:       attackHitlagFactor[c.NormalCounter][i],
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: true,
		}
		ap := combat.NewCircleHitOnTargetFanAngle(
			c.Core.Combat.Player(),
			geometry.Point{Y: attackOffsets[c.NormalCounter][i]},
			attackHitboxes[c.NormalCounter][0],
			attackFanAngles[c.NormalCounter][i],
		)
		if c.NormalCounter == 1 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: attackOffsets[c.NormalCounter][i]},
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
