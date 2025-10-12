package ineffa

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
	attackHitmarks        = [][]int{{9}, {17}, {10, 10 + 11}, {27}}
	attackHitlagHaltFrame = [][]float64{{0.06}, {0.06}, {0.00, 0.06}, {0.06}}
	attackDefHalt         = [][]bool{{true}, {true}, {false, true}, {true}}
	attackHitboxes        = [][][]float64{{{1.6, 2}}, {{2}}, {{2.5}, {2.5}}, {{2}}}
	attackOffsets         = [][]float64{{1}, {-0.2}, {-0.2, -0.2}, {-0.2}}
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 35)
	attackFrames[0][action.ActionAttack] = 12
	attackFrames[0][action.ActionCharge] = 27

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 46)
	attackFrames[1][action.ActionAttack] = 19
	attackFrames[1][action.ActionCharge] = 39

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 48)
	attackFrames[2][action.ActionAttack] = 33
	attackFrames[2][action.ActionCharge] = 39

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 55)
	attackFrames[3][action.ActionAttack] = 52
}

// Normal attack damage queue generator
// relatively standard with no major differences versus other characters
func (c *char) Attack(p map[string]int) (action.Info, error) {
	for i, mult := range attack[c.NormalCounter] {
		ai := info.AttackInfo{
			ActorIndex:         c.Index(),
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

		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			info.Point{Y: attackOffsets[c.NormalCounter][i]},
			attackHitboxes[c.NormalCounter][i][0],
		)
		if c.NormalCounter == 0 {
			ai.StrikeType = attacks.StrikeTypeSpear
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				info.Point{Y: attackOffsets[c.NormalCounter][i]},
				attackHitboxes[c.NormalCounter][i][0],
				attackHitboxes[c.NormalCounter][i][1],
			)
		}
		c.QueueCharTask(func() {
			c.Core.QueueAttack(ai, ap, 0, 0)
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
