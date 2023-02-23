package albedo

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
	attackHitmarks        = []int{12, 11, 17, 17, 27}
	attackHitlagHaltFrame = []float64{0.03, 0.03, 0.06, 0.09, 0.12}
	attackHitboxes        = [][]float64{{1.6}, {2}, {3, 2.6}, {1.5, 3.7}, {2}}
	attackOffsets         = []float64{0.8, 0.1, -0.2, 0.1, 0.8}
	attackFanAngles       = []float64{360, 180, 360, 360, 360}
)

const normalHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 14) // N1 -> N2
	attackFrames[0][action.ActionCharge] = 23                             // N1 -> CA

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 22) // N2 -> N3
	attackFrames[1][action.ActionCharge] = 21                             // N2 -> CA

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 32) // N3 -> N4
	attackFrames[2][action.ActionCharge] = 41                             // N3 -> CA

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 34) // N4 -> N5
	attackFrames[3][action.ActionCharge] = 36                             // N4 -> CA

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[3], 62) // N5 -> N1
	attackFrames[4][action.ActionCharge] = 500                            //TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:          attacks.AttackTagNormal,
		ICDTag:             attacks.ICDTagNormalAttack,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeSlash,
		Element:            attributes.Physical,
		Durability:         25,
		Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter] * 60,
		CanBeDefenseHalted: true,
	}
	ap := combat.NewCircleHitOnTargetFanAngle(
		c.Core.Combat.Player(),
		combat.Point{Y: attackOffsets[c.NormalCounter]},
		attackHitboxes[c.NormalCounter][0],
		attackFanAngles[c.NormalCounter],
	)
	if c.NormalCounter == 2 || c.NormalCounter == 3 {
		ap = combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			combat.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
			attackHitboxes[c.NormalCounter][1],
		)
	}
	//we don't need to use char queue here since each hit is single hit
	c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter], attackHitmarks[c.NormalCounter])

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}
