package lisa

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackHitmark = []int{26, 18, 17, 31}
var attackReleaseFrame = []int{15, 12, 17, 31}
var attackFrames [][]int

const normalHitNum = 4

//TODO: these frames are wrong!!!
func init() {
	attackFrames = make([][]int, normalHitNum)
	attackFrames[0] = frames.InitNormalCancelSlice(attackReleaseFrame[0], 30)
	attackFrames[1] = frames.InitNormalCancelSlice(attackReleaseFrame[1], 20)
	attackFrames[2] = frames.InitNormalCancelSlice(attackReleaseFrame[2], 34)
	attackFrames[3] = frames.InitNormalCancelSlice(attackReleaseFrame[3], 57)
}

func (c *char) Attack(p map[string]int) action.ActionInfo {

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  combat.AttackTagNormal,
		ICDTag:     combat.ICDTagLisaElectro,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	//todo: Does it really snapshot immediately?
	c.Core.QueueAttack(ai,
		combat.NewDefSingleTarget(1, combat.TargettableEnemy),
		0,
		attackHitmark[c.NormalCounter],
	)

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackReleaseFrame[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}
