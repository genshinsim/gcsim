package qiqi

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = [][]int{{11}, {18}, {41, 42}, {39, 40}, {29}}
var attackHitlagHaltFrame = [][]float64{{0.03}, {0.03}, {0.03, 0.03}, {0.03, 0.03}, {0.12}}
var attackHitlagFactor = [][]float64{{0.01}, {0.01}, {0.05, 0.05}, {0.05, 0.05}, {0.01}}

const normalHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 11)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 18)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 42)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][1], 40)
	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 29)
}

// Standard attack - nothing special
func (c *char) Attack(p map[string]int) action.ActionInfo {

	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			Mult:               mult[c.TalentLvlAttack()],
			AttackTag:          combat.AttackTagNormal,
			ICDTag:             combat.ICDTagNormalAttack,
			ICDGroup:           combat.ICDGroupDefault,
			StrikeType:         combat.StrikeTypeSlash,
			Element:            attributes.Physical,
			Durability:         25,
			HitlagFactor:       attackHitlagFactor[c.NormalCounter][i],
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: true,
		}
		c.QueueCharTask(func() {
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHit(c.Core.Combat.Player(), 0.3, false, combat.TargettableEnemy),
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
