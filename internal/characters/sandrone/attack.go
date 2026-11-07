package sandrone

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
	attackHitmarks        = []int{24, 25, 49}
	attackPoiseDMG        = []float64{102, 89.9, 137.5}
	attackHitlagHaltFrame = []float64{0.1, 0.1, .12}
	attackHitboxes        = [][]float64{{4.3, 2.5}, {2.7}, {2.5}}
	attackOffsets         = []float64{0.5, 0.15, 2}
)

const normalHitNum = 3

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 53) // N1 -> CA
	attackFrames[0][action.ActionAttack] = 30
	attackFrames[0][action.ActionWalk] = 38

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 46) // N2 -> CA
	attackFrames[1][action.ActionAttack] = 31
	attackFrames[1][action.ActionWalk] = 46

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 87) // N3 -> CA
	attackFrames[2][action.ActionAttack] = 69
	attackFrames[2][action.ActionWalk] = 79
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex:         c.Index(),
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
		info.Point{Y: attackOffsets[c.NormalCounter]},
		attackHitboxes[c.NormalCounter][0],
	)
	if c.NormalCounter == 0 {
		ap = combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			info.Point{Y: attackOffsets[c.NormalCounter]},
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
