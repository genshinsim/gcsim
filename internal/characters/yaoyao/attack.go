package yaoyao

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	attackFrames          [][]int
	attackHitmarks        = [][]int{{13}, {16}, {12, 12 + 19}, {21}}
	attackHitlagHaltFrame = [][]float64{{0.03}, {0.03}, {0, 0.03}, {0.03}}
	attackHitboxes        = [][][]float64{
		{{1.2, 3}},
		{{1.6, 3.3}},
		{{1.6, 3.3}, {1.2, 3.3}},
		{{1.6, 3.3}},
	}
)

const (
	normalHitNum = 4
)

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 22)
	attackFrames[0][action.ActionAttack] = 22
	attackFrames[0][action.ActionCharge] = 22
	attackFrames[0][action.ActionWalk] = 28

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 22)
	attackFrames[1][action.ActionAttack] = 22
	attackFrames[1][action.ActionCharge] = 25
	attackFrames[1][action.ActionWalk] = 31

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 28)
	attackFrames[2][action.ActionAttack] = 41
	attackFrames[2][action.ActionCharge] = 37
	attackFrames[2][action.ActionWalk] = 51

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 37)
	attackFrames[3][action.ActionAttack] = 59
	attackFrames[3][action.ActionAttack] = 54
	attackFrames[3][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			Mult:               mult[c.TalentLvlAttack()],
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeSpear,
			Element:            attributes.Physical,
			Durability:         25,
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: true,
		}
		ap := combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			nil,
			attackHitboxes[c.NormalCounter][i][0],
			attackHitboxes[c.NormalCounter][i][1],
		)
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
