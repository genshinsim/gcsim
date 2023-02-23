package ayaka

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
	attackHitmarks        = [][]int{{8}, {10}, {16}, {8, 15, 22}, {27}}
	attackHitlagHaltFrame = [][]float64{{0.03}, {0.03}, {0.06}, {0, 0, 0.03}, {0}}
	attackHitlagFactor    = [][]float64{{0.01}, {0.01}, {0.01}, {0, 0, 0.05}, {0.01}}
	attackDefHalt         = [][]bool{{true}, {true}, {true}, {false, false, true}, {false}}
	attackRadius          = []float64{1.6, 1.2, 2.8, 1.6, 0.8}
	attackOffsets         = []float64{0.8, 0.8, -0.5, 0.6, 0}
	attackFanAngles       = []float64{360, 360, 60, 360, 360}
)

const normalHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 22)
	attackFrames[0][action.ActionAttack] = 9

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 20)
	attackFrames[1][action.ActionAttack] = 19

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 32)
	attackFrames[2][action.ActionCharge] = 31

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][2], 23)
	attackFrames[3][action.ActionAttack] = 22

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 66)
	attackFrames[4][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	for i, mult := range attack[c.NormalCounter] {
		icdGroup := combat.ICDGroupDefault
		centerTarget := c.Core.Combat.Player()
		if c.NormalCounter == 4 {
			icdGroup = combat.ICDGroupPoleExtraAttack    // N5 has a different ICDGroup
			centerTarget = c.Core.Combat.PrimaryTarget() // N5 is a bullet
		}
		ai := combat.AttackInfo{
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			ActorIndex:         c.Index,
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             combat.ICDTagNormalAttack,
			ICDGroup:           icdGroup,
			StrikeType:         attacks.StrikeTypeSlash,
			Element:            attributes.Physical,
			Durability:         25,
			HitlagFactor:       attackHitlagFactor[c.NormalCounter][i],
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
			Mult:               mult[c.TalentLvlAttack()],
		}
		ap := combat.NewCircleHitFanAngle(
			c.Core.Combat.Player(),
			centerTarget,
			combat.Point{Y: attackOffsets[c.NormalCounter]},
			attackRadius[c.NormalCounter],
			attackFanAngles[c.NormalCounter],
		)
		c.QueueCharTask(func() {
			c.Core.QueueAttack(ai, ap, 0, 0, c.c1)
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
