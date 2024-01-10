package chevreuse

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
	attackFrames          [][]int
	attackHitmarks        = [][]int{{20}, {14}, {16, 23}, {15}} // TODO
	attackHitlagHaltFrame = [][]float64{{0.06}, {0.06}, {0, 0.06}, {0.1}}
	attackDefHalt         = [][]bool{{true}, {true}, {false, true}, {true}}
	attackHitboxes        = [][]float64{{1.6, 2}, {1.6}, {1.6, 2}, {1.5, 4}}
	attackStrikeTypes     = [][]attacks.StrikeType{
		{attacks.StrikeTypeSlash},
		{attacks.StrikeTypeSlash},
		{attacks.StrikeTypeSlash, attacks.StrikeTypeSlash},
		{attacks.StrikeTypeSpear}}
	attackOffsets = [][]float64{{0}, {0.4}, {0, 0.4}, {1.0}}
)

func init() {
	// NA cancels (polearm)
	attackFrames = make([][]int, normalHitNum)
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 32) // N1 -> CA
	attackFrames[0][action.ActionAttack] = 20

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 33) // N2 -> N3/CA

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 48) // N3 -> N4
	attackFrames[2][action.ActionCharge] = 43

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 69) // N4 -> N1
	attackFrames[3][action.ActionCharge] = 500
	// TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attackStrikeTypes[c.NormalCounter][i],
			PoiseDMG:           0.25,
			Element:            attributes.Physical,
			Durability:         25,
			Mult:               mult[c.TalentLvlAttack()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}

		var ap combat.AttackPattern

		if c.NormalCounter == 0 || (c.NormalCounter == 2 && i == 0) || c.NormalCounter == 3 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: attackOffsets[c.NormalCounter][i]},
				attackHitboxes[c.NormalCounter][0],
				attackHitboxes[c.NormalCounter][1],
			)
		}

		if c.NormalCounter == 1 || (c.NormalCounter == 2 && i == 1) {
			ap = combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: attackOffsets[c.NormalCounter][i]},
				attackHitboxes[c.NormalCounter][0],
			)
		}
		c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter][i], attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}, nil
}
