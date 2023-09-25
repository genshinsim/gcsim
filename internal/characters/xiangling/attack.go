package xiangling

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
	attackHitmarks        = [][]int{{12}, {8}, {11, 18}, {5, 15, 24, 29}, {21}}
	attackHitlagHaltFrame = [][]float64{{0.03}, {0.03}, {0.03, 0}, {0, 0, 0, 0.03}, {0.09}}
	attackHitboxes        = [][][]float64{
		{{1.2, 3}},
		{{1.6, 3.3}},
		{{1.6, 3.3}, {1.2, 3.3}},
		{{1.2, 3.3}, {1.2, 3.3}, {1.2, 3.3}, {1.2, 3.3}},
		{{1.6, 3.3}},
	}
)

const (
	normalHitNum = 5
	c2Debuff     = "xiangling-c2"
)

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 20)

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 17)

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 28)
	attackFrames[2][action.ActionCharge] = 24

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][3], 37)
	attackFrames[3][action.ActionCharge] = 34

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 70)
	attackFrames[4][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	done := false
	var c2CB func(a combat.AttackCB)
	if c.Base.Cons >= 2 && c.NormalCounter == 4 {
		c2CB = c.c2(done)
	}
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
			c.Core.QueueAttack(ai, ap, 0, 0, c2CB)
		}, attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}, nil
}
