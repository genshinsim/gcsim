package durin

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
	attackHitmarks        = [][]int{{11}, {9}, {14, 14 + 23}, {38}}
	attackHitlagHaltFrame = [][]float64{{0.02}, {0.03}, {0.00, 0.03}, {0.05}}
	attackHitboxes        = [][][]float64{{{2.15}}, {{2.15}}, {{2.4}, {2.4}}, {{3.3, 4.8}}}
	attackOffsets         = [][]info.Point{{{Y: 0.3}}, {{Y: 0.6}}, {{X: -0.2, Y: 0.4}, {X: -0.2, Y: 0.4}}, {{X: -0.2, Y: -0.6}}}
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 29)
	attackFrames[0][action.ActionAttack] = 19
	attackFrames[0][action.ActionCharge] = 21

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 30)
	attackFrames[1][action.ActionAttack] = 14
	attackFrames[1][action.ActionCharge] = 21

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 55)
	attackFrames[2][action.ActionAttack] = 48
	attackFrames[2][action.ActionCharge] = 43

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 66)
	attackFrames[3][action.ActionWalk] = 65
}

// Normal attack damage queue generator
// relatively standard with no major differences versus other characters
func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillWindowKey) {
		return c.skillRecastBlack(), nil
	}

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
			CanBeDefenseHalted: false,
		}

		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			attackOffsets[c.NormalCounter][i],
			attackHitboxes[c.NormalCounter][i][0],
		)
		if c.NormalCounter == 3 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				attackOffsets[c.NormalCounter][i],
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
