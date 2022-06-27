package zhongli

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = [][]int{{11}, {9}, {8}, {16}, {11, 18, 23, 29}, {29}}

const normalHitNum = 6

func init() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 30)
	attackFrames[0][action.ActionAttack] = 18

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 30)
	attackFrames[1][action.ActionAttack] = 13

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 28)
	attackFrames[2][action.ActionAttack] = 19

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 34)
	attackFrames[3][action.ActionCharge] = 33

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][3], 31)
	attackFrames[4][action.ActionAttack] = 27

	attackFrames[5] = frames.InitNormalCancelSlice(attackHitmarks[5][0], 54)
	attackFrames[5][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  combat.AttackTagNormal,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
		FlatDmg:    0.0139 * c.MaxHP(),
	}

	for i := 0; i < hits[c.NormalCounter]; i++ {
		c.Core.QueueAttack(
			ai,
			combat.NewDefCircHit(0.1, false, combat.TargettableEnemy),
			attackHitmarks[c.NormalCounter][i],
			attackHitmarks[c.NormalCounter][i],
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
