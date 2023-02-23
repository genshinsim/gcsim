package candace

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
	attackHitmarks        = [][]int{{11}, {16}, {16, 39}, {43}}
	attackHitlagHaltFrame = [][]float64{{0.03}, {0.03}, {0, 0.03}, {0.04}}
	attackHitlagDefHalt   = [][]bool{{true}, {true}, {false, true}, {true}}
	attackStrikeTypes     = [][]attacks.StrikeType{
		{attacks.StrikeTypeSlash},
		{attacks.StrikeTypeBlunt},
		{attacks.StrikeTypeSlash, attacks.StrikeTypeSlash},
		{attacks.StrikeTypeSpear},
	}
	attackHitboxes  = [][]float64{{2.5}, {2.2, 3}, {2.5}, {2.2, 7}}
	attackOffsets   = []float64{-0.3, -0.1, 0.3, -0.1}
	attackFanAngles = []float64{360, 360, 270, 360}
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 32) // N1 -> CA
	attackFrames[0][action.ActionAttack] = 20

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 33) // N2 -> N3/CA

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 48) // N3 -> N4
	attackFrames[2][action.ActionCharge] = 43

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 69) // N4 -> N1
	attackFrames[3][action.ActionCharge] = 500                               // TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           combat.ICDGroupDefault,
			StrikeType:         attackStrikeTypes[c.NormalCounter][i],
			Element:            attributes.Physical,
			Durability:         25,
			Mult:               mult[c.TalentLvlAttack()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackHitlagDefHalt[c.NormalCounter][i],
		}
		ap := combat.NewCircleHitOnTargetFanAngle(
			c.Core.Combat.Player(),
			combat.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
			attackFanAngles[c.NormalCounter],
		)
		if c.NormalCounter == 1 || c.NormalCounter == 3 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				combat.Point{Y: attackOffsets[c.NormalCounter]},
				attackHitboxes[c.NormalCounter][0],
				attackHitboxes[c.NormalCounter][1],
			)
		}
		c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter][i], attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}
