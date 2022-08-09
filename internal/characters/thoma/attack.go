package thoma

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	attackFrames          [][]int
	attackHitmarks        = [][]int{{13}, {18}, {10, 13}, {13}}
	attackHitlagHaltFrame = [][]float64{{0.06}, {0.09}, {0, 0}, {0}}
	attackDefHalt         = [][]bool{{true}, {true}, {false, false}, {false}}
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	// N1 -> x
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 21)
	attackFrames[0][action.ActionAttack] = 20
	attackFrames[0][action.ActionCharge] = 21

	// N2 -> x
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 27)
	attackFrames[1][action.ActionAttack] = 27
	attackFrames[1][action.ActionCharge] = 25

	// N3 -> x
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 32)
	attackFrames[2][action.ActionAttack] = 31
	attackFrames[2][action.ActionCharge] = 32

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 58)
	attackFrames[3][action.ActionAttack] = 58
	attackFrames[3][action.ActionCharge] = 500 // TODO: this action is illegal; need better way to handle it
}

// Normal attack damage queue generator
// relatively standard with no major differences versus other characters
func (c *char) Attack(p map[string]int) action.ActionInfo {
	lastMultiHit := 0
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
		hitmark := lastMultiHit + attackHitmarks[c.NormalCounter][i]
		c.QueueCharTask(func() {
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHit(c.Core.Combat.Player(), 0.1, false, combat.TargettableEnemy),
				0,
				0,
			)
		}, hitmark)
		lastMultiHit = hitmark
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}
