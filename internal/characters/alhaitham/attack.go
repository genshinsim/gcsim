package alhaitham

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = [][]int{{8}, {11}, {14, 28}, {18}, {34}}
var attackHitlagHaltFrame = [][]float64{{.03}, {.03}, {0, .03}, {.06}, {0}}
var attackDefHalt = [][]bool{{true}, {true}, {false, true}, {true}, {false}}
var attackHitboxes = [][]float64{{2.5}, {2, 2.5}, {2, 3}, {3, 4.5}, {2.5}}
var attackOffsets = []float64{-0.1, -0.1, 0, 0, 2.2}
var attackFanAngles = []float64{180, 360, 360, 360, 360}

const normalHitNum = 5

func init() {
	// TODO: Attack frames and cancels, koli pls
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 25)
	attackFrames[0][action.ActionAttack] = 16
	attackFrames[0][action.ActionCharge] = 25

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 21)
	attackFrames[1][action.ActionAttack] = 21

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 47)
	attackFrames[2][action.ActionAttack] = 47

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 28)
	attackFrames[3][action.ActionAttack] = 28

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 34)
	attackFrames[4][action.ActionAttack] = 34

}

func (c *char) Attack(p map[string]int) action.ActionInfo {

	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:          combat.AttackTagNormal,
			ICDTag:             combat.ICDTagNormalAttack,
			ICDGroup:           combat.ICDGroupDefault,
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
		if c.NormalCounter == 1 || c.NormalCounter == 2 || c.NormalCounter == 3 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				combat.Point{Y: attackOffsets[c.NormalCounter]},
				attackHitboxes[c.NormalCounter][0],
				attackHitboxes[c.NormalCounter][1],
			)
		}
		if c.NormalCounter == 4 {
			ap = combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				combat.Point{Y: attackOffsets[c.NormalCounter]},
				attackHitboxes[c.NormalCounter][0])
		}
		c.QueueCharTask(func() {
			c.Core.QueueAttack(
				ai,
				ap,
				0,
				0,
				c.projectionAttack,
				c.a1CB,
			)
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
