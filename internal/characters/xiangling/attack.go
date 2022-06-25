package xiangling

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = [][]int{{12}, {8}, {11, 18}, {5, 15, 24, 29}, {21}}

const normalHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 20)

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 17)

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 28)
	attackFrames[2][action.ActionCharge] = 24

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][3], 37)
	attackFrames[3][action.ActionCharge] = 34

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 70)
	attackFrames[4][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
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

	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.QueueAttack(
			ai,
			combat.NewDefCircHit(0.1, false, combat.TargettableEnemy),
			attackHitmarks[c.NormalCounter][i],
			attackHitmarks[c.NormalCounter][i],
		)
	}

	//if n = 5, add explosion for c2
	if c.Base.Cons >= 2 && c.NormalCounter == 4 {
		// According to TCL, does not snapshot and has no ability type scaling tags
		// TODO: Does not mention ICD or pyro aura strength?
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Oil Meets Fire (C2)",
			AttackTag:  combat.AttackTagNone,
			ICDTag:     combat.ICDTagNormalAttack,
			ICDGroup:   combat.ICDGroupDefault,
			Element:    attributes.Pyro,
			Durability: 25,
			Mult:       .75,
		}
		c.Core.QueueAttack(ai, combat.NewDefCircHit(0.1, false, combat.TargettableEnemy), 120, 120) //todo: explosion frames
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}
