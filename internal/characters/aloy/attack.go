package aloy

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	attackFrames   [][]int
	attackHitmarks = [][]int{{11, 24}, {16}, {23}, {30}}
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][1], 31) // N1 -> N2
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 28) // N2 -> N3
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 38) // N3 -> N4
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 61) // N4 -> N1

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
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypePierce,
		Element:    attributes.Physical,
		Durability: 25,
	}

	if c.StatusIsActive(rushingIceKey) {
		ai.ICDTag = combat.ICDTagNormalAttack
	}

	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.QueueAttack(
			ai,
			combat.NewBoxHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				nil,
				0.1,
				1,
			),
			attackHitmarks[c.NormalCounter][i],
			attackHitmarks[c.NormalCounter][i]+travel)
	}

	defer c.AdvanceNormalIndex()

	// return animation cd
	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}
