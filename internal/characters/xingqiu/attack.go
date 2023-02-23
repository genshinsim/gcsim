package xingqiu

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	attackFrames           [][]int
	attackHitmarks         = [][]int{{10}, {13}, {9, 19}, {17}, {18, 39}}
	attackHitlagHaltFrames = []float64{0.03, 0.03, 0.06, 0.06, 0.1}
	attackHitboxes         = [][][]float64{{{1.5}}, {{1.5}}, {{1.5}, {1.5}}, {{1, 2}}, {{1, 2}, {2}}}
	attackOffsets          = [][]float64{{0.8}, {0.8}, {0.6, 0.6}, {0}, {0, 0.8}}
)

const normalHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 35)
	attackFrames[0][action.ActionAttack] = 18

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 29)
	attackFrames[1][action.ActionAttack] = 24

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 35)
	attackFrames[2][action.ActionCharge] = 26

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 33)
	attackFrames[3][action.ActionAttack] = 28

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][1], 66)
	attackFrames[4][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		AttackTag:          attacks.AttackTagNormal,
		ICDTag:             attacks.ICDTagNormalAttack,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeSlash,
		Element:            attributes.Physical,
		Durability:         25,
		HitlagHaltFrames:   attackHitlagHaltFrames[c.NormalCounter] * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}

	for i, mult := range attack[c.NormalCounter] {
		ax := ai
		ax.Abil = fmt.Sprintf("Normal %v", c.NormalCounter)
		ax.Mult = mult[c.TalentLvlAttack()]
		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			combat.Point{Y: attackOffsets[c.NormalCounter][i]},
			attackHitboxes[c.NormalCounter][i][0],
		)
		if c.NormalCounter == 3 || (c.NormalCounter == 4 && i == 0) {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				combat.Point{Y: attackOffsets[c.NormalCounter][i]},
				attackHitboxes[c.NormalCounter][i][0],
				attackHitboxes[c.NormalCounter][i][1],
			)
		}
		c.QueueCharTask(func() {
			c.Core.QueueAttack(ax, ap, 0, 0)
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
