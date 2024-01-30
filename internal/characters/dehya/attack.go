package dehya

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
	attackHitmarks        = []int{22, 26, 26, 41}
	attackPoiseDMG        = []float64{93.33, 92.72, 115.14, 143.17}
	attackHitlagHaltFrame = []float64{0.09, 0.10, 0.08, .12}
	attackHitboxes        = [][]float64{{2.2}, {2.3, 4.3}, {1.8}, {3, 4.3}}
	attackOffsets         = []float64{0.5, -1.3, 0.5, -0.8}
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 31) // N1 -> N2
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 34) // N2 -> N3
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 43) // N3 -> N4
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 85) // N4 -> N1
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	burstAction := c.UseBurstAction()
	if burstAction != nil {
		return *burstAction, nil
	}
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:          attacks.AttackTagNormal,
		ICDTag:             attacks.ICDTagNormalAttack,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		PoiseDMG:           attackPoiseDMG[c.NormalCounter],
		Element:            attributes.Physical,
		Durability:         25,
		Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter] * 60,
		CanBeDefenseHalted: true,
	}
	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		geometry.Point{Y: attackOffsets[c.NormalCounter]},
		attackHitboxes[c.NormalCounter][0],
	)
	if c.NormalCounter == 1 || c.NormalCounter == 3 {
		ap = combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
			attackHitboxes[c.NormalCounter][1],
		)
	}
	c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter], attackHitmarks[c.NormalCounter])

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}, nil
}

// TODO: charged attack
