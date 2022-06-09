package razor

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var hitmarks = []int{25, 46, 38, 83}

func (c *char) attackFrameFunc(next action.Action) int {
	// back out what last attack was
	n := c.NormalCounter - 1
	if n < 0 {
		n = c.NormalHitNum - 1
	}
	return frames.AtkSpdAdjust(
		attackFrames[n][next],
		c.Stat(attributes.AtkSpd),
	)
}

func (c *char) initNormalCancels() {
	// normal cancels
	attackFrames = make([][]int, c.NormalHitNum) // should be 4

	// n1 animations
	attackFrames[0] = frames.InitNormalCancelSlice(hitmarks[0], 25)
	// n2 animations
	attackFrames[1] = frames.InitNormalCancelSlice(hitmarks[1], 46) // 71-25
	// n3 animations
	attackFrames[2] = frames.InitNormalCancelSlice(hitmarks[2], 38) // 109-71
	// n4 animations
	attackFrames[3] = frames.InitNormalCancelSlice(hitmarks[3], 83) // 192-109
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
		hitmarks[c.NormalCounter],
		hitmarks[c.NormalCounter],
	)

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          c.attackFrameFunc,
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   hitmarks[c.NormalCounter],
		Post:            hitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}
