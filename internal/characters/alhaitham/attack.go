package alhaitham

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
	attackFrames   [][]int
	attackHitmarks = [][]int{{11}, {13}, {15, 29}, {21}, {35}}

	attackHitlagHaltFrame = [][]float64{{0.02}, {0.02}, {0, 0.02}, {0}, {0.02}}
	attackDefHalt         = [][]bool{{true}, {true}, {false, true}, {false}, {true}}
	attackHitboxes        = [][]float64{{2.5}, {2, 2.5}, {2, 3}, {3, 4.5}, {2.5}}

	attackOffsets   = []float64{-0.1, -0.1, 0, 0, 2.2}
	attackFanAngles = []float64{180, 360, 360, 360, 360}
)

const normalHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 29)
	attackFrames[0][action.ActionAttack] = 15
	attackFrames[0][action.ActionCharge] = 23

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 31)
	attackFrames[1][action.ActionAttack] = 22
	attackFrames[1][action.ActionCharge] = 25

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 55)
	attackFrames[2][action.ActionAttack] = 44
	attackFrames[2][action.ActionCharge] = 45

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 51)
	attackFrames[3][action.ActionAttack] = 30
	attackFrames[3][action.ActionCharge] = 37

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 90)
	attackFrames[4][action.ActionAttack] = 67
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	strikeType := attacks.StrikeTypeSlash
	if c.NormalCounter == 1 {
		strikeType = attacks.StrikeTypeSpear
	}
	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         strikeType,
			Element:            attributes.Physical,
			Durability:         25,
			Mult:               mult[c.TalentLvlAttack()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}

		ap := combat.NewCircleHitOnTargetFanAngle(
			c.Core.Combat.Player(),
			geometry.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
			attackFanAngles[c.NormalCounter],
		)

		if c.NormalCounter == 1 || c.NormalCounter == 2 || c.NormalCounter == 3 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: attackOffsets[c.NormalCounter]},
				attackHitboxes[c.NormalCounter][0],
				attackHitboxes[c.NormalCounter][1],
			)
		}

		c.Core.QueueAttack(
			ai,
			ap,
			attackHitmarks[c.NormalCounter][i],
			attackHitmarks[c.NormalCounter][i],
			c.projectionAttack,
		)
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}
