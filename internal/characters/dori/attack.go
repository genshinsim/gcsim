package dori

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = []int{12, 13, 13}
var attackHitlagHaltFrame = []float64{0.03, 0.03, 0.06}

const normalHitNum = 3

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 23) // N1 -> CA
	attackFrames[0][action.ActionAttack] = 19                             // N1 -> N2

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 27) // N2 -> CA
	attackFrames[1][action.ActionAttack] = 17                             // N2 -> N3

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 50) // N3 -> CA
	attackFrames[2][action.ActionAttack] = 42                             // N3 -> N4
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	for _, mult := range auto[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:       c.Index,
			Abil:             fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:        combat.AttackTagNormal,
			ICDTag:           combat.ICDTagNormalAttack,
			ICDGroup:         combat.ICDGroupDefault,
			StrikeType:       combat.StrikeTypeBlunt,
			Element:          attributes.Physical,
			Durability:       25,
			Mult:             mult[c.TalentLvlAttack()],
			HitlagFactor:     0.01,
			HitlagHaltFrames: attackHitlagHaltFrame[c.NormalCounter] * 60,
		}
		//c6 key check
		if c.StatusIsActive(c6key) {
			ai.Element = attributes.Electro
		}
		c.QueueCharTask(func() {
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHit(c.Core.Combat.Player(), 1, false, combat.TargettableEnemy),
				0,
				0,
			)
		}, attackHitmarks[c.NormalCounter])
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}
