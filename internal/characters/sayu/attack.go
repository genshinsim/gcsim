package sayu

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = []int{23, 70 - 23, 109 - 70, 187 - 109}

//TODO: need frame count
func init() {
	attackFrames = make([][]int, normalHitNum)
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 23)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 70-23)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 109-70)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 187-109)
}

func (c *char) Attack(p map[string]int) action.ActionInfo {

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  combat.AttackTagNormal,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Physical,
		Durability: 25,
	}
	snap := c.Snapshot(&ai)
	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.QueueAttackWithSnap(ai,
			snap,
			combat.NewDefCircHit(0.3, false, combat.TargettableEnemy),
			attackHitmarks[c.NormalCounter]-2+i,
		)
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}
