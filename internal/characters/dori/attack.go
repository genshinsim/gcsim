package dori

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	attackFrames          [][]int
	attackHitmarks        = [][]int{{27}, {19, 33}, {60}}
	attackHitlagHaltFrame = [][]float64{{0.1}, {0, 0}, {0.08}}
	attackHitlagDefHalt   = [][]bool{{true}, {false, false}, {true}}
	attackOffsets         = []float64{0.5, 0.5, 1}
)

const normalHitNum = 3

func init() {
	attackFrames = make([][]int, normalHitNum)
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 44)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][1], 46)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 108)
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	c6CB := c.makeC6CB()
	for i, mult := range auto[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeBlunt,
			Element:            attributes.Physical,
			Durability:         25,
			Mult:               mult[c.TalentLvlAttack()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackHitlagDefHalt[c.NormalCounter][i],
		}
		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			combat.Point{Y: attackOffsets[c.NormalCounter]},
			2,
		)
		c.QueueCharTask(func() {
			c.Core.QueueAttack(ai, ap, 0, 0, c6CB)
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
