package odette

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	attackFrames          [][]int
	attackHitmarks        = [][]int{{9}, {9}, {16, 16 + 15}, {15}, {15}}
	attackHitlagHaltFrame = [][]float64{{0.03}, {0.03}, {0, 0}, {0.06}, {0.08}}
	attackDefHalt         = [][]bool{{true}, {true}, {false, false}, {true}, {true}}
	attackHitboxes        = [][]float64{{1.7}, {1.7}, {1.6, 2.8}, {2, 2.6}, {6, 2}}
	attackOffsets         = []float64{0.6, 0.8, 0.3, -0.2, 0.6}
)

const (
	normalHitNum      = 5
	shunsuikenHitmark = 5
)

func init() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 25) // N1 -> CA
	attackFrames[0][action.ActionAttack] = 17                                // N1 -> N2
	attackFrames[0][action.ActionWalk] = 19                                  // N1 -> W

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 23) // N2 -> W
	attackFrames[1][action.ActionAttack] = 19                                // N2 -> N3
	attackFrames[1][action.ActionCharge] = 21                                // N2 -> CA

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 69) // N3 -> CA
	attackFrames[2][action.ActionAttack] = 54                                // N3 -> N4
	attackFrames[2][action.ActionWalk] = 63                                  // N3 -> W

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 85) // N4 -> W
	attackFrames[3][action.ActionAttack] = 53                                // N4 -> N5
	attackFrames[3][action.ActionCharge] = 73                                // N4 -> CA

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 70) // N5 -> W
	attackFrames[4][action.ActionAttack] = 62                                // N5 -> N1
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	for i, mult := range attack[c.NormalCounter] {
		ai := info.AttackInfo{
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			ActorIndex:         c.Index(),
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
			info.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
		)
		if c.NormalCounter >= 2 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				info.Point{Y: attackOffsets[c.NormalCounter]},
				attackHitboxes[c.NormalCounter][0],
				attackHitboxes[c.NormalCounter][1],
			)
		}
		c.QueueCharTask(func() {
			c.Core.QueueAttack(ai, ap, 0, 0)
		}, attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	// normal state
	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}, nil
}
