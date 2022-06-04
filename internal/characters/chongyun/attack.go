package chongyun

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = []int{24, 38, 62, 80}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		ActorIndex: c.Index,
		AttackTag:  combat.AttackTagNormal,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewDefCircHit(1, false, combat.TargettableEnemy),
		attackHitmarks[c.NormalCounter],
		attackHitmarks[c.NormalCounter],
	)

	if c.Base.Cons >= 1 && c.NormalCounter == 3 {
		ai := combat.AttackInfo{
			Abil:       "Chongyun C1",
			ActorIndex: c.Index,
			AttackTag:  combat.AttackTagNone,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       .5,
		}
		//3 blades
		for i := 0; i < 3; i++ {
			c.Core.QueueAttack(
				ai,
				combat.NewDefCircHit(1, false, combat.TargettableEnemy),
				attackHitmarks[c.NormalCounter]+i*5,
				attackHitmarks[c.NormalCounter]+i*5,
			)
		}
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		Post:            attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}
