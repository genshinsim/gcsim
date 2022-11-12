package yunjin

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = [][]int{{15}, {13}, {8, 23}, {11, 23}, {15}}
var attackHitlagHaltFrame = [][]float64{{0.03}, {0.03}, {0, 0.03}, {0, 0.03}, {0.04}}
var attackDefHalt = [][]bool{{true}, {true}, {false, true}, {false, false}, {true}}

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
			AttackTag:          combat.AttackTagNormal,
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
		c.QueueCharTask(func() {
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHit(c.Core.Combat.Player(), 0.1),
				0,
				0,
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
