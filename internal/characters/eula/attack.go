package eula

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	attackFrames          [][]int
	attackHitmarks        = [][]int{{30}, {19}, {25, 42}, {17}, {29, 56}}
	attackHitlagHaltFrame = [][]float64{{0.09}, {.12}, {0, .09}, {.09}, {0, .12}}
	attackDefHalt         = [][]bool{{true}, {true}, {false, true}, {true}, {false, true}}
	attackRadius          = []float64{2, 2, 2, 2, 1.8}
)

const normalHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 34)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 36)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 56)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 44)
	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][1], 105)
}

func (c *char) Attack(p map[string]int) action.ActionInfo {

	for i, mult := range auto[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:          combat.AttackTagNormal,
			ICDTag:             combat.ICDTagNormalAttack,
			ICDGroup:           combat.ICDGroupDefault,
			StrikeType:         combat.StrikeTypeBlunt,
			Element:            attributes.Physical,
			Durability:         25,
			Mult:               mult[c.TalentLvlAttack()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}
		radius := attackRadius[c.NormalCounter]
		c.QueueCharTask(func() {
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHit(c.Core.Combat.Player(), radius),
				0,
				0,
			)
		}, attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}
