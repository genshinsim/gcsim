package kirara

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
	attackHitmarks        = [][]int{{13}, {18}, {17, 22}, {39}}
	attackHitlagHaltFrame = [][]float64{{0.03}, {0.03}, {0.01, 0.05}, {0.06}}
	attackHitboxes        = [][]float64{{2}, {1.8, 3.8}, {2}, {2.1}}
	attackOffsets         = [][]float64{{0.6}, {-0.3}, {0.4, 0.9}, {1}}
)

const normalHitNum = 4

func init() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 34) // N1 -> Walk
	attackFrames[0][action.ActionAttack] = 20
	attackFrames[0][action.ActionCharge] = 24

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 35) // N2 -> Walk
	attackFrames[1][action.ActionAttack] = 29
	attackFrames[1][action.ActionCharge] = 28

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 61) // N3 -> Walk
	attackFrames[2][action.ActionAttack] = 54
	attackFrames[2][action.ActionCharge] = 40

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 63) // N4 -> Walk
	attackFrames[3][action.ActionAttack] = 60
	attackFrames[3][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:       c.Index,
			Abil:             fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:        attacks.AttackTagNormal,
			ICDTag:           attacks.ICDTagNormalAttack,
			ICDGroup:         attacks.ICDGroupDefault,
			StrikeType:       attacks.StrikeTypeSlash,
			Element:          attributes.Physical,
			Durability:       25,
			Mult:             mult[c.TalentLvlAttack()],
			HitlagFactor:     0.01,
			HitlagHaltFrames: attackHitlagHaltFrame[c.NormalCounter][i] * 60,
		}
		ap := combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.Player(),
			geometry.Point{Y: attackOffsets[c.NormalCounter][i]},
			attackHitboxes[c.NormalCounter][0],
		)
		if c.NormalCounter == 2 {
			ap = combat.NewBoxHit(
				c.Core.Combat.Player(),
				c.Core.Combat.Player(),
				geometry.Point{Y: attackOffsets[c.NormalCounter][i]},
				attackHitboxes[c.NormalCounter][0],
				attackHitboxes[c.NormalCounter][1],
			)
		}
		c.QueueCharTask(func() {
			c.Core.QueueAttack(ai, ap, 0, 0)
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
