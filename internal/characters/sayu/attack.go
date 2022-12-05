package sayu

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	attackFrames          [][]int
	attackHitmarks        = [][]int{{23}, {29}, {14, 26}, {35}}
	attackHitlagHaltFrame = [][]float64{{0.1}, {0.1}, {0, 0.08}, {0.08}}
	attackDefHalt         = [][]bool{{true}, {true}, {false, true}, {true}}
	attackRadius          = [][]float64{{2.03}, {2}, {2.15, 2}, {2.65}}
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 36) // N1 -> N2
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 48) // N2 -> N3
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 52) // N3 -> N4
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 71) // N4 -> N1
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	for i, mult := range attack[c.NormalCounter] {
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
		radius := attackRadius[c.NormalCounter][i]
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
