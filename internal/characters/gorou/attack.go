package gorou

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = []int{17, 12, 27, 31}

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 21) // N1 -> N2
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 26) // N2 -> N3
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 43) // N3 -> N4
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 55) // N4 -> N1
}

// Normal attack damage queue generator
// relatively standard with no major differences versus other bow characters
// Has "travel" parameter, used to set the number of frames that the arrow is in the air (default = 10)
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
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.PrimaryTarget(), 0.5),
		attackHitmarks[c.NormalCounter],
		attackHitmarks[c.NormalCounter]+travel,
	)

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}
