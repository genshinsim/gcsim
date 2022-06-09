package sucrose

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = []int{17, 18, 28, 28}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  combat.AttackTagNormal,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewDefSingleTarget(1, combat.TargettableEnemy),
		attackHitmarks[c.NormalCounter],
		attackHitmarks[c.NormalCounter],
	)

	defer c.AdvanceNormalIndex()

	if c.Base.Cons >= 4 {
		c.c4()
	}

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		Post:            attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}
