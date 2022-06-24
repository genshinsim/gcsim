package kazuha

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = [][]int{{12}, {11}, {16, 25}, {15}, {15, 23, 31}}

const normalHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 16)
	attackFrames[0][action.ActionAttack] = 15

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 25)
	attackFrames[1][action.ActionAttack] = 20

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 35)
	attackFrames[2][action.ActionAttack] = 30

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 40)
	attackFrames[3][action.ActionCharge] = 35

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][2], 71)
	attackFrames[4][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  combat.AttackTagNormal,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeSlash,
		Element:    attributes.Physical,
		Durability: 25,
	}

	act := action.ActionInfo{
		Frames:              frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength:     attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:       attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:               action.NormalAttackState,
		FramePausedOnHitlag: c.FramePausedOnHitlag,
	}

	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.QueueAttack(
			ai,
			combat.NewDefCircHit(0.3, false, combat.TargettableEnemy),
			attackHitmarks[c.NormalCounter][i],
			attackHitmarks[c.NormalCounter][i],
		)
	}

	defer c.AdvanceNormalIndex()

	return act
}
