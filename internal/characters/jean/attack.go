package jean

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
	attackHitmarks        = []int{13, 6, 17, 37, 25}
	attackHitlagHaltFrame = []float64{.03, .03, .06, .06, .1}
	attackRadius          = []float64{1.5, 2.2, 2.8, 1.6, 1.6}
	attackOffsets         = []float64{1.5, -0.5, -1, 0.6, 0.6}
	attackFanAngles       = []float64{360, 150, 30, 360, 360}
)

const normalHitNum = 5

func init() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 25)
	attackFrames[0][action.ActionAttack] = 22

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 20)
	attackFrames[1][action.ActionAttack] = 14

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 31)
	attackFrames[2][action.ActionAttack] = 28

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 49)
	attackFrames[3][action.ActionAttack] = 44

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4], 68)
	attackFrames[4][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:          attacks.AttackTagNormal,
		ICDTag:             attacks.ICDTagNormalAttack,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeSlash,
		Element:            attributes.Physical,
		Durability:         25,
		Mult:               auto[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter] * 60,
		CanBeDefenseHalted: true,
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTargetFanAngle(
			c.Core.Combat.Player(),
			geometry.Point{Y: attackOffsets[c.NormalCounter]},
			attackRadius[c.NormalCounter],
			attackFanAngles[c.NormalCounter],
		),
		attackHitmarks[c.NormalCounter],
		attackHitmarks[c.NormalCounter],
		c.makeA1CB(),
	)

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}
