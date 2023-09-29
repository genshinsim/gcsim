package xiao

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
	attackHitmarks        = [][]int{{4, 17}, {15}, {15}, {14, 31}, {16}, {39}}
	attackHitlagHaltFrame = [][]float64{{0, 0.01}, {0.01}, {0.01}, {0.02, 0.02}, {0.02}, {0.04}}
	attackDefHalt         = [][]bool{{false, true}, {true}, {true}, {false, true}, {true}, {true}}
	attackStrikeTypes     = [][]attacks.StrikeType{
		{attacks.StrikeTypeSlash, attacks.StrikeTypeSpear},
		{attacks.StrikeTypeSlash},
		{attacks.StrikeTypeSlash},
		{attacks.StrikeTypeSlash, attacks.StrikeTypeSlash},
		{attacks.StrikeTypeSpear},
		{attacks.StrikeTypeSlash},
	}
	attackHitboxes = [][][][]float64{
		{
			{{1.8}, {1.4, 2.7}},
			{{1.6}},
			{{1.6}},
			{{1.6}, {1.8}},
			{{1.5, 3}},
			{{2}},
		},
		{
			{{2}, {1.6, 3}},
			{{1.8}},
			{{1.8}},
			{{1.8}, {2}},
			{{1.7, 3.2}},
			{{2.4}},
		},
	}
	attackOffsets = [][][]float64{
		{{0, -0.1}, {0, 0}},
		{{0, 1}},
		{{0, 1.1}},
		{{-0.1, 0.9}, {0, 0.8}},
		{{0, 0}},
		{{0, 1.1}},
	}
	attackFanAngles = [][]float64{{150, 360}, {300}, {300}, {320, 320}, {360}, {360}}
)

const normalHitNum = 6

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][1], 26)
	attackFrames[0][action.ActionAttack] = 25

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 27)
	attackFrames[1][action.ActionAttack] = 22

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 38)
	attackFrames[2][action.ActionAttack] = 26

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][1], 42)
	attackFrames[3][action.ActionAttack] = 39

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 30)
	attackFrames[4][action.ActionAttack] = 24

	attackFrames[5] = frames.InitNormalCancelSlice(attackHitmarks[5][0], 79)
	attackFrames[5][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
}

// Normal attack damage queue generator
// relatively standard with no major differences versus other characters
func (c *char) Attack(p map[string]int) (action.Info, error) {
	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attackStrikeTypes[c.NormalCounter][i],
			Element:            attributes.Physical,
			Durability:         25,
			Mult:               mult[c.TalentLvlAttack()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}

		burstIndex := 0
		if c.StatusIsActive(burstBuffKey) {
			burstIndex = 1
		}
		ap := combat.NewCircleHitOnTargetFanAngle(
			c.Core.Combat.Player(),
			geometry.Point{X: attackOffsets[c.NormalCounter][i][0], Y: attackOffsets[c.NormalCounter][i][1]},
			attackHitboxes[burstIndex][c.NormalCounter][i][0],
			attackFanAngles[c.NormalCounter][i],
		)
		if (c.NormalCounter == 0 && i == 1) || c.NormalCounter == 4 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{X: attackOffsets[c.NormalCounter][i][0], Y: attackOffsets[c.NormalCounter][i][1]},
				attackHitboxes[burstIndex][c.NormalCounter][i][0],
				attackHitboxes[burstIndex][c.NormalCounter][i][1],
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
