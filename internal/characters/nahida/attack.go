package nahida

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

const normalHitNum = 4

var attackFrames [][]int
var attackHitmarks = []int{2 * (1091 - 1079), 2 * (1105 - 1099), 2 * (1123 - 1113), 2 * (1149 - 1128)}

func init() {
	attackFrames = make([][]int, normalHitNum)

	//TODO: update frames
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 20*2)

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 14*2)

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 15*2)

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[2], 36*2)
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  combat.AttackTagNormal,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewDefSingleTarget(c.Core.Combat.DefaultTarget),
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
