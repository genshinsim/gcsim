package diona

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = []int{16, 37 - 16, 67 - 37, 101 - 67, 152 - 101}

const normalHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 16)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 37-16)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 67-37)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 101-67)
	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4], 152-101)
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  combat.AttackTagNormal,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypePierce,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       auto[c.NormalCounter][c.TalentLvlAttack()],
	}
	a := action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
	c.Core.QueueAttack(
		ai,
		combat.NewDefSingleTarget(1, combat.TargettableEnemy),
		attackHitmarks[c.NormalCounter],
		travel+attackHitmarks[c.NormalCounter],
	)

	defer c.AdvanceNormalIndex()

	return a
}
