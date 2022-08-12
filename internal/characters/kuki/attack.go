package kuki

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = []int{12, 13, 13, 23}
var attackHitlagHaltFrame = []float64{0.03, 0.03, 0.06, 0.10}

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 23) // N1 -> CA
	attackFrames[0][action.ActionAttack] = 19                             // N1 -> N2

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 27) // N2 -> CA
	attackFrames[1][action.ActionAttack] = 17                             // N2 -> N3

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 50) // N3 -> CA
	attackFrames[2][action.ActionAttack] = 42                             // N3 -> N4

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 57) // N4 -> N1
	attackFrames[3][action.ActionCharge] = 500                            // N4 -> CA, TODO: this action is illegal; need better way to handle it
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
		Mult:               auto[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter] * 60,
		CanBeDefenseHalted: true,
	}
	// no multihits so no need for char queue here
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), .3, false, combat.TargettableEnemy),
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
