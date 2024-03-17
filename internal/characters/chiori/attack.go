package chiori

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
	attackHitmarks        = [][]int{{17}, {16}, {24, 32}, {37}}
	attackHitlagHaltFrame = [][]float64{{0.03}, {0.03}, {0.06, 0}, {0}}
	attackHitlagFactor    = [][]float64{{0.01}, {0.01}, {0.01, 0.01}, {0.05}}
	attackDefHalt         = [][]bool{{true}, {true}, {true, false}, {true}}
	attackHitboxes        = [][]float64{{1.6}, {1.8}, {1.8, 3.6}, {2.8}}
	attackOffsets         = []float64{0.5, 0.6, -0.6, 0.9}
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 22) // N1 -> CA
	attackFrames[0][action.ActionAttack] = 19                                // N1 -> N2 rerecorded

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 33) // N2 -> CA
	attackFrames[1][action.ActionAttack] = 22                                // N2 -> N3 rerecorded

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 42) // N3 -> CA
	attackFrames[2][action.ActionAttack] = 41                                // N3 -> N4

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 59) // N4 -> N1
	attackFrames[3][action.ActionCharge] = 500                               // TODO: this action is illegal; need better way to handle it
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
			HitlagFactor:       attackHitlagFactor[c.NormalCounter][i],
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}

		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
		)
		if c.NormalCounter == 2 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: attackOffsets[c.NormalCounter]},
				attackHitboxes[c.NormalCounter][0],
				attackHitboxes[c.NormalCounter][1],
			)
		}

		c.Core.Tasks.Add(func() {
			snap := c.Snapshot(&ai)
			c.c6NAIncrease(&ai, &snap)
			c.Core.QueueAttackWithSnap(ai, snap, ap, 0)
		}, attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	c.tryTriggerA1TailoringNA()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}, nil
}
