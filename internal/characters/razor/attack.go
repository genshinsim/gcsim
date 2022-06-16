package razor

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = []int{25, 46, 38, 83}

const normalHitNum = 4

func (c *char) init() {
	attackFrames = make([][]int, normalHitNum) // should be 4

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 25)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 46) // 71-25
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 38) // 109-71
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 83) // 192-109
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
	}

	ai.Mult = attack[c.NormalCounter][c.TalentLvlAttack()]
	c.Core.QueueAttack(
		ai,
		combat.NewDefCircHit(0.5, false, combat.TargettableEnemy),
		attackHitmarks[c.NormalCounter],
		attackHitmarks[c.NormalCounter],
	)

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}
