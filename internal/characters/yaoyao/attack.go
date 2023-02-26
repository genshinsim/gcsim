package yaoyao

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
	attackHitmarks        = [][]int{{13}, {16}, {12, 12 + 19}, {21}}
	attackHitlagHaltFrame = [][]float64{{0.06}, {0.06}, {0, 0.01}, {0.09}}
	attackDefHalt         = [][]bool{{true}, {true}, {false, true}, {true}}
	attackHitboxes        = [][]float64{{1.2, 3}, {2}, {2}, {2.2}}
	attackOffsets         = []float64{0, 0.5, 0.5, 0.5, 1.5}
	attackFanAngles       = []float64{360, 270, 270, 360}
)

const (
	normalHitNum = 4
)

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 28)
	attackFrames[0][action.ActionAttack] = 22
	attackFrames[0][action.ActionCharge] = 22

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 31)
	attackFrames[1][action.ActionAttack] = 22
	attackFrames[1][action.ActionCharge] = 25

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 51)
	attackFrames[2][action.ActionAttack] = 41
	attackFrames[2][action.ActionCharge] = 37

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 59)
	attackFrames[3][action.ActionWalk] = 54
	attackFrames[3][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
}

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
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}
		ap := combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][i],
			attackFanAngles[c.NormalCounter],
		)
		if c.NormalCounter == 0 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: attackOffsets[c.NormalCounter]},
				attackHitboxes[c.NormalCounter][0],
				attackHitboxes[c.NormalCounter][1],
			)
		}

		c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter][i], attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}
