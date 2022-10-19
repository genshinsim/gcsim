package candace

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	attackFrames          [][]int
	attackHitmarks        = [][]int{{11}, {16}, {16, 39}, {43}}
	attackHitlagHaltFrame = [][]float64{{0.03}, {0.03}, {0, 0.03}, {0.04}}
	attackHitlagDefHalt   = [][]bool{{true}, {true}, {false, true}, {true}}
	attackStrikeTypes     = [][]combat.StrikeType{
		{combat.StrikeTypeSlash},
		{combat.StrikeTypeBlunt},
		{combat.StrikeTypeSlash, combat.StrikeTypeSlash},
		{combat.StrikeTypeSpear},
	}
	// {{radius 2.5 circle}, {x=2.2 z=3.0 box}, {radius 2.5 fanAngle 270 circle, radius 2.5 fanAngle 270 circle}, {x=2.2 z=7.0 box}}
	attackRadius = []float64{2.5, 1.5, 2.5, 3.5}
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
			AttackTag:          combat.AttackTagNormal,
			ICDTag:             combat.ICDTagNormalAttack,
			ICDGroup:           combat.ICDGroupDefault,
			StrikeType:         attackStrikeTypes[c.NormalCounter][i],
			Element:            attributes.Physical,
			Durability:         25,
			Mult:               mult[c.TalentLvlAttack()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackHitlagDefHalt[c.NormalCounter][i],
		}
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(
				c.Core.Combat.Player(),
				attackRadius[c.NormalCounter],
				false,
				combat.TargettableEnemy,
				combat.TargettableGadget,
			),
			attackHitmarks[c.NormalCounter][i],
			attackHitmarks[c.NormalCounter][i],
		)
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}
