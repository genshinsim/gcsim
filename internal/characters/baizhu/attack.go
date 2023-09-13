package baizhu

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

const normalHitNum = 4

var (
	attackFrames   [][]int
	attackHitmarks = [][]int{{17}, {25}, {23, 35}, {30}}
	attackHitboxes = [][]float64{{2, 3}, {2, 3}, {2.4, 3.0}, {3.2, 3.0}}
	attackOffsets  = []float64{-0.2, -0.2, -0.2, -0.2}
)

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 33) // N1 -> Walk
	attackFrames[0][action.ActionAttack] = 22
	attackFrames[0][action.ActionCharge] = 19

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 39) // N2 -> Walk
	attackFrames[1][action.ActionAttack] = 33
	attackFrames[1][action.ActionCharge] = 31

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 46) // N3 -> Walk
	attackFrames[2][action.ActionAttack] = 40
	attackFrames[2][action.ActionCharge] = 39

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 67) // N4 -> N1
	attackFrames[3][action.ActionWalk] = 62
	attackFrames[3][action.ActionCharge] = 57
}

func (c *char) Attack(p map[string]int) action.Info {
	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:   c.Index,
			Abil:         fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:    attacks.AttackTagNormal,
			ICDTag:       attacks.ICDTagNormalAttack,
			ICDGroup:     attacks.ICDGroupDefault,
			StrikeType:   attacks.StrikeTypeDefault,
			Element:      attributes.Dendro,
			Durability:   25,
			Mult:         mult[c.TalentLvlAttack()],
			HitlagFactor: 0.01,
		}

		ap := combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
			attackHitboxes[c.NormalCounter][1],
		)
		c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter][i], attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}
