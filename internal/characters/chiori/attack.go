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
	attackFrames   [][]int
	attackHitmarks = [][]int{{25 - 1}, {20 - 1}, {20 - 1, 59 - 20 - 1}, {62 - 1}} //TODO:

	//TODO: cp from alhaitham
	attackHitlagHaltFrame = [][]float64{{0.03}, {0.03}, {0.06, 0}, {0}}
	attackDefHalt         = [][]bool{{true}, {true}, {false, true}, {false}}
	attackHitboxes        = [][]float64{{2.5}, {2}, {2, 3}, {3}}
	attackOffsets         = []float64{-0.1, -0.1, 0, 0}
	attackFanAngles       = []float64{180, 360, 360, 360}
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 25)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 20)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 59)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 62)
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

		//TODO: wtf is this even
		ap := combat.NewCircleHitOnTargetFanAngle(
			c.Core.Combat.Player(),
			geometry.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
			attackFanAngles[c.NormalCounter],
		)

		c.Core.Tasks.Add(func() {
			snap := c.Snapshot(&ai)
			c.c6AttackMod(&ai, &snap)
			c.Core.QueueAttackWithSnap(ai, snap, ap, 0)
		}, attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	c.tryTriggerA1Tailoring()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}, nil
}
