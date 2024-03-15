package gaming

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
	attackHitmarks        = []int{25, 18, 34, 48}
	attackHitlagHaltFrame = []float64{0.10, 0.09, 0.12, 0.12}
	attackPoiseDMG        = []float64{112.13300323486328, 105.6989974975586, 142.60800170898438, 171.09500122070312}
	attackHitboxes        = [][]float64{{2}, {2}, {2.5, 4}, {2}}
	attackOffsets         = []float64{1, 1, -1.7, 2}
	attackFanAngles       = []float64{360, 270, 360, 360}
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 39)
	attackFrames[0][action.ActionAttack] = 30

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 46)
	attackFrames[1][action.ActionAttack] = 32

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 86)
	attackFrames[2][action.ActionAttack] = 79

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 89)
	attackFrames[3][action.ActionAttack] = 87
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		ActorIndex:         c.Index,
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
	ap := combat.NewCircleHitOnTargetFanAngle(
		c.Core.Combat.Player(),
		geometry.Point{Y: attackOffsets[c.NormalCounter]},
		attackHitboxes[c.NormalCounter][0],
		attackFanAngles[c.NormalCounter],
	)
	if c.NormalCounter == 2 {
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
