package travelergeo

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
	attackFrames          [][][]int
	attackHitmarks        = [][]int{{13, 13, 16, 30, 25}, {16, 10, 19, 23, 14}}
	attackHitlagHaltFrame = [][]float64{{0.03, 0.03, 0.06, 0.09, 0.12}, {0.03, 0.03, 0.06, 0.06, 0.10}}
	a4Hitmark             = []int{18, 20}
	attackHitboxes        = [][][]float64{{{1.4, 2.2}, {1.7}, {1.5, 2.2}, {1.7}, {1.75}}, {{1.6}, {1.4, 2.2}, {1.5}, {1.5}, {1.6}}}
	attackOffsets         = [][]float64{{0, 0.6, 0.4, 0.6, 0.6}, {1, 0, 0.7, 0.7, 1}}
	attackFanAngles       = [][]float64{{360, 180, 360, 360, 240}, {360, 360, 360, 360, 360}}
)

const normalHitNum = 5

func init() {
	attackFrames = make([][][]int, 2)

	// Male
	attackFrames[0] = make([][]int, normalHitNum)

	attackFrames[0][0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 28) // N1 -> CA
	attackFrames[0][0][action.ActionAttack] = 17                                // N1 -> N2

	attackFrames[0][1] = frames.InitNormalCancelSlice(attackHitmarks[0][1], 28) // N2 -> CA
	attackFrames[0][1][action.ActionAttack] = 26                                // N2 -> N3

	attackFrames[0][2] = frames.InitNormalCancelSlice(attackHitmarks[0][2], 36) // N3 -> CA
	attackFrames[0][2][action.ActionAttack] = 32                                // N3 -> N4

	attackFrames[0][3] = frames.InitNormalCancelSlice(attackHitmarks[0][3], 45) // N4 -> CA
	attackFrames[0][3][action.ActionAttack] = 39                                // N4 -> N5

	attackFrames[0][4] = frames.InitNormalCancelSlice(attackHitmarks[0][4], 69) // N5 -> N1
	attackFrames[0][4][action.ActionCharge] = 500                               // N5 -> CA, TODO: this action is illegal; need better way to handle it

	// Female
	attackFrames[1] = make([][]int, normalHitNum)

	attackFrames[1][0] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 32) // N1 -> CA
	attackFrames[1][0][action.ActionAttack] = 24                                // N1 -> N2

	attackFrames[1][1] = frames.InitNormalCancelSlice(attackHitmarks[1][1], 23) // N2 -> CA
	attackFrames[1][1][action.ActionAttack] = 21                                // N2 -> N3

	attackFrames[1][2] = frames.InitNormalCancelSlice(attackHitmarks[1][2], 39) // N3 -> CA
	attackFrames[1][2][action.ActionAttack] = 27                                // N3 -> N4

	attackFrames[1][3] = frames.InitNormalCancelSlice(attackHitmarks[1][3], 45) // N4 -> CA
	attackFrames[1][3][action.ActionAttack] = 38                                // N4 -> N5

	attackFrames[1][4] = frames.InitNormalCancelSlice(attackHitmarks[1][4], 64) // N5 -> N1
	attackFrames[1][4][action.ActionCharge] = 500                               // N5 -> CA, TODO: this action is illegal; need better way to handle it
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
		Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   attackHitlagHaltFrame[c.gender][c.NormalCounter] * 60,
		CanBeDefenseHalted: true,
	}
	ap := combat.NewCircleHitOnTargetFanAngle(
		c.Core.Combat.Player(),
		geometry.Point{Y: attackOffsets[c.gender][c.NormalCounter]},
		attackHitboxes[c.gender][c.NormalCounter][0],
		attackFanAngles[c.gender][c.NormalCounter],
	)
	if (c.gender == 0 && (c.NormalCounter == 0 || c.NormalCounter == 2)) ||
		(c.gender == 1 && c.NormalCounter == 1) {
		ap = combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: attackOffsets[c.gender][c.NormalCounter]},
			attackHitboxes[c.gender][c.NormalCounter][0],
			attackHitboxes[c.gender][c.NormalCounter][1],
		)
	}
	c.Core.QueueAttack(
		ai,
		ap,
		attackHitmarks[c.gender][c.NormalCounter],
		attackHitmarks[c.gender][c.NormalCounter],
	)
	c.a4()

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames[c.gender]),
		AnimationLength: attackFrames[c.gender][c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.gender][c.NormalCounter],
		State:           action.NormalAttackState,
	}
}
