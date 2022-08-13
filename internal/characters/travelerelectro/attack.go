package travelerelectro

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][][]int
var attackHitmarks = [][]int{{13, 13, 16, 30, 25}, {16, 10, 19, 23, 14}}
var attackHitlagHaltFrame = [][]float64{{0.03, 0.03, 0.06, 0.09, 0.12}, {0.03, 0.03, 0.06, 0.06, 0.10}}

const normalHitNum = 5

func init() {
	attackFrames = make([][][]int, 2)

	// Male
	attackFrames[0] = make([][]int, normalHitNum)

	attackFrames[0][0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 28) // N1 -> CA
	attackFrames[0][0][action.ActionAttack] = 17                                // N1 -> N2

	attackFrames[0][1] = frames.InitNormalCancelSlice(attackHitmarks[0][1], 28) // N2 -> CA
	attackFrames[0][1][action.ActionAttack] = 26                                // N2 -> N3

	attackFrames[0][2] = frames.InitNormalCancelSlice(attackHitmarks[0][2], 36) // N3 -> CA
	attackFrames[0][2][action.ActionAttack] = 32                                // N3 -> N4

	attackFrames[0][3] = frames.InitNormalCancelSlice(attackHitmarks[0][3], 45) // N4 -> CA
	attackFrames[0][3][action.ActionAttack] = 39                                // N4 -> N5

	attackFrames[0][4] = frames.InitNormalCancelSlice(attackHitmarks[0][4], 69) // N5 -> N1
	attackFrames[0][4][action.ActionCharge] = 500                               // N5 -> CA, TODO: this action is illegal; need better way to handle it

	// Female
	attackFrames[1] = make([][]int, normalHitNum)

	attackFrames[1][0] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 32) // N1 -> CA
	attackFrames[1][0][action.ActionAttack] = 24                                // N1 -> N2

	attackFrames[1][1] = frames.InitNormalCancelSlice(attackHitmarks[1][1], 23) // N2 -> CA
	attackFrames[1][1][action.ActionAttack] = 21                                // N2 -> N3

	attackFrames[1][2] = frames.InitNormalCancelSlice(attackHitmarks[1][2], 39) // N3 -> CA
	attackFrames[1][2][action.ActionAttack] = 27                                // N3 -> N4

	attackFrames[1][3] = frames.InitNormalCancelSlice(attackHitmarks[1][3], 45) // N4 -> CA
	attackFrames[1][3][action.ActionAttack] = 38                                // N4 -> N5

	attackFrames[1][4] = frames.InitNormalCancelSlice(attackHitmarks[1][4], 64) // N5 -> N1
	attackFrames[1][4][action.ActionCharge] = 500                               // N5 -> CA, TODO: this action is illegal; need better way to handle it
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
		Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   attackHitlagHaltFrame[c.female][c.NormalCounter] * 60,
		CanBeDefenseHalted: true,
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), 0.3, false, combat.TargettableEnemy),
		attackHitmarks[c.female][c.NormalCounter],
		attackHitmarks[c.female][c.NormalCounter],
	)

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames[c.female]),
		AnimationLength: attackFrames[c.female][c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.female][c.NormalCounter],
		State:           action.NormalAttackState,
	}
}
