package rosaria

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
	attackHitmarks        = [][]int{{9}, {13}, {19, 28}, {32}, {26, 40}}
	attackHitlagHaltFrame = [][]float64{{0.06}, {0.06}, {0, 0.03}, {0.09}, {0.06, 0.06}}
	attackDefHalt         = [][]bool{{true}, {true}, {false, true}, {true}, {false, true}}
	attackHitboxes        = [][][]float64{{{2.5}}, {{2.5}}, {{2.2, 3.6}, {2.2, 3.6}}, {{2.3}}, {{2, 3}, {2.5}}}
	attackOffsets         = [][]float64{{0.5}, {0.5}, {0, 0}, {2}, {0, 1.5}}
	attackFanAngles       = []float64{270, 270, 360, 360, 360}
)

const normalHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 24)
	attackFrames[0][action.ActionAttack] = 19

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 27)
	attackFrames[1][action.ActionAttack] = 23

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 34)
	attackFrames[2][action.ActionAttack] = 31

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 52)
	attackFrames[3][action.ActionAttack] = 44

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][1], 66)
	attackFrames[4][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
}

// Normal attack damage queue generator
// relatively standard with no major differences versus other characters
func (c *char) Attack(p map[string]int) action.Info {
	c1CB := c.makeC1CB()
	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			Mult:               mult[c.TalentLvlAttack()],
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeSlash,
			Element:            attributes.Physical,
			Durability:         25,
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}
		if c.NormalCounter == 2 {
			ai.StrikeType = attacks.StrikeTypeSpear
		}
		ap := combat.NewCircleHitOnTargetFanAngle(
			c.Core.Combat.Player(),
			geometry.Point{Y: attackOffsets[c.NormalCounter][i]},
			attackHitboxes[c.NormalCounter][i][0],
			attackFanAngles[c.NormalCounter],
		)
		if c.NormalCounter == 2 || (c.NormalCounter == 4 && i == 0) {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: attackOffsets[c.NormalCounter][i]},
				attackHitboxes[c.NormalCounter][i][0],
				attackHitboxes[c.NormalCounter][i][1],
			)
		}
		c.QueueCharTask(func() {
			c.Core.QueueAttack(ai, ap, 0, 0, c1CB)
		}, attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}
