package emilie

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
	attackHitmarks = []int{17, 15, 15, 39}

	attackOffsets = [][]float64{
		{0, 0.4},
		{0, 0},
		{0, 0.4},
		{0, 0.4},
	}
	attackHitboxes = [][]float64{
		{2},
		{1.6, 3}, // box
		{2.5},
		{3},
	}

	attackHitlagHaltFrame = []float64{0.06, 0.06, 0.06, 0.09}
)

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[1], 26)
	attackFrames[0][action.ActionAttack] = 25

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[0], 27)
	attackFrames[1][action.ActionAttack] = 22

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[0], 38)
	attackFrames[2][action.ActionAttack] = 26

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[1], 42)
	attackFrames[3][action.ActionAttack] = 39
	attackFrames[3][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		ActorIndex:         c.Index,
		AttackTag:          attacks.AttackTagNormal,
		ICDTag:             attacks.ICDTagNormalAttack,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeSlash,
		Element:            attributes.Physical,
		Durability:         25,
		Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter] * 60,
		CanBeDefenseHalted: true,
	}
	c.applyC6Bonus(&ai)

	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.PrimaryTarget(),
		geometry.Point{X: attackOffsets[c.NormalCounter][0], Y: attackOffsets[c.NormalCounter][1]},
		attackHitboxes[c.NormalCounter][0],
	)
	if c.NormalCounter == 1 {
		ai.StrikeType = attacks.StrikeTypeSpear
		ap = combat.NewBoxHitOnTarget(
			c.Core.Combat.PrimaryTarget(),
			geometry.Point{X: attackOffsets[c.NormalCounter][0], Y: attackOffsets[c.NormalCounter][1]},
			attackHitboxes[c.NormalCounter][0],
			attackHitboxes[c.NormalCounter][1],
		)
	}

	c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter], attackHitmarks[c.NormalCounter])

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[len(attackHitmarks)-1],
		State:           action.NormalAttackState,
	}, nil
}
