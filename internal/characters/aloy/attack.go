package aloy

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

//TODO: not sure where these are from; no idea how accurate
var attackHitmarks = []int{30, 18, 37, 43}
var attackFrames [][]int

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)
	//TODO: no cancelled frames here; this is machine gunning
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], attackHitmarks[0])
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], attackHitmarks[1])
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], attackHitmarks[2])
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], attackHitmarks[3])
}

// Standard attack - infusion mechanics are handled as part of the skill
func (c *char) Attack(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  combat.AttackTagNormal,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Physical,
		Durability: 25,
	}

	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.QueueAttack(
			ai,
			combat.NewDefSingleTarget(c.Core.Combat.DefaultTarget, combat.TargettableEnemy),
			attackHitmarks[c.NormalCounter]+i,
			attackHitmarks[c.NormalCounter]+i+travel)
	}

	defer c.AdvanceNormalIndex()

	// return animation cd
	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}
