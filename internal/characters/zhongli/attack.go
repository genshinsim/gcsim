package zhongli

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	attackFrames          [][]int
	attackEarliestCancel  = []int{11, 9, 8, 16, 4, 29}
	attackHitmarks        = [][]int{{11}, {9}, {8}, {16}, {11, 18, 23, 29}, {29}}
	attackHitlagHaltFrame = [][]float64{{0.02}, {0.02}, {0.02}, {0.02}, {0, 0, 0, 0}, {0.02}}
	attackDefHalt         = [][]bool{{true}, {true}, {true}, {true}, {false, false, false, false}, {true}}
	attackRadius          = []float64{2.04, 2, 0.9, 1.7, 2.06, 2.06}
)

const normalHitNum = 6

func init() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackEarliestCancel[0], 30)
	attackFrames[0][action.ActionAttack] = 18

	attackFrames[1] = frames.InitNormalCancelSlice(attackEarliestCancel[1], 30)
	attackFrames[1][action.ActionAttack] = 13

	attackFrames[2] = frames.InitNormalCancelSlice(attackEarliestCancel[2], 28)
	attackFrames[2][action.ActionAttack] = 19

	attackFrames[3] = frames.InitNormalCancelSlice(attackEarliestCancel[3], 34)
	attackFrames[3][action.ActionCharge] = 33

	attackFrames[4] = frames.InitNormalCancelSlice(attackEarliestCancel[4], 31)
	attackFrames[4][action.ActionAttack] = 27
	attackFrames[4][action.ActionSkill] = 5
	attackFrames[4][action.ActionBurst] = 5
	attackFrames[4][action.ActionDash] = 5
	attackFrames[4][action.ActionJump] = 5

	attackFrames[5] = frames.InitNormalCancelSlice(attackEarliestCancel[5], 54)
	attackFrames[5][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	for i := 0; i < hits[c.NormalCounter]; i++ {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:          combat.AttackTagNormal,
			ICDTag:             combat.ICDTagNormalAttack,
			ICDGroup:           combat.ICDGroupDefault,
			StrikeType:         combat.StrikeTypeSpear,
			Element:            attributes.Physical,
			Durability:         25,
			Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
			FlatDmg:            0.0139 * c.MaxHP(),
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}
		if c.NormalCounter == 1 || c.NormalCounter == 4 {
			ai.StrikeType = combat.StrikeTypeSlash
		}
		radius := attackRadius[c.NormalCounter]
		//the multihit part generates no hitlag so this is fine
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), radius),
			attackHitmarks[c.NormalCounter][i],
			attackHitmarks[c.NormalCounter][i],
		)
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackEarliestCancel[c.NormalCounter],
		State:           action.NormalAttackState,
	}

}
