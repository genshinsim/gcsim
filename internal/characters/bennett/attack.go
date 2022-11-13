package bennett

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = []int{13, 9, 13, 25, 24}
var attackHitlagHaltFrames = []float64{0.03, 0.03, 0.06, 0.09, 0.12}

const normalHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 33)
	attackFrames[0][action.ActionAttack] = 20

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 27)
	attackFrames[1][action.ActionAttack] = 17

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 46)
	attackFrames[2][action.ActionAttack] = 37

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 48)
	attackFrames[3][action.ActionAttack] = 44

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4], 60)
	attackFrames[4][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:          combat.AttackTagNormal,
		ICDTag:             combat.ICDTagNormalAttack,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeSlash,
		Element:            attributes.Physical,
		Durability:         25,
		HitlagHaltFrames:   attackHitlagHaltFrames[c.NormalCounter] * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
		Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), 0.1),
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
