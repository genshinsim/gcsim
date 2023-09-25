package lynette

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var (
	attackFrames          [][]int
	attackHitmarks        = [][]int{{17}, {13}, {23, 31}, {29}}
	attackHitlagHaltFrame = [][]float64{{0.03}, {0.03}, {0.06, 0}, {0}}
	attackDefHalt         = [][]bool{{true}, {true}, {true, false}, {false}}
	attackHitboxes        = [][]float64{{2}, {1.8, 2.8}, {1.5, 2.5}, {2.2, 2.5}}
	attackOffsets         = []float64{0.3, -0.3, 0, 0}
)

const normalHitNum = 4

func init() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 30) // N1 -> Walk
	attackFrames[0][action.ActionAttack] = 24
	attackFrames[0][action.ActionCharge] = 25

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 27) // N2 -> Walk
	attackFrames[1][action.ActionAttack] = 16
	attackFrames[1][action.ActionCharge] = 25

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 47) // N3 -> CA/Walk
	attackFrames[2][action.ActionAttack] = 35

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 68) // N4 -> Walk
	attackFrames[3][action.ActionAttack] = 60
	attackFrames[3][action.ActionCharge] = 500 // TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeSlash,
			Element:            attributes.Physical,
			Durability:         25,
			Mult:               mult[c.TalentLvlAttack()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}
		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
		)
		if c.NormalCounter >= 1 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: attackOffsets[c.NormalCounter]},
				attackHitboxes[c.NormalCounter][0],
				attackHitboxes[c.NormalCounter][1],
			)
		}
		c.QueueCharTask(func() {
			c.Core.QueueAttack(ai, ap, 0, 0)
		}, attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}, nil
}
