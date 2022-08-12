package shenhe

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = [][]int{{14}, {17}, {19}, {14, 18}, {26}}
var attackHitlagHaltFrame = [][]float64{{0.02}, {0.02}, {0.02}, {0, 0.02}, {0.1}}
var attackDefHalt = [][]bool{{true}, {true}, {true}, {false, true}, {true}}

const normalHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 23)
	attackFrames[0][action.ActionCharge] = 22

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 29)
	attackFrames[1][action.ActionAttack] = 17

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 33)
	attackFrames[2][action.ActionAttack] = 32

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][1], 27)
	attackFrames[3][action.ActionAttack] = 25

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 49)
	attackFrames[4][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
}

// Normal attack damage queue generator
// relatively standard with no major differences versus other characters
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
		c.QueueCharTask(func() {
			c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 0.1, false, combat.TargettableEnemy), 0, 0)
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
