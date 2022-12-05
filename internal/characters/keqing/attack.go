package keqing

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	attackFrames          [][]int
	attackHitmarks        = [][]int{{11}, {11}, {15}, {12, 22}, {26}}
	attackHitlagHaltFrame = [][]float64{{.03}, {.03}, {.06}, {0, .03}, {0}}
	attackDefHalt         = [][]bool{{true}, {true}, {true}, {false, true}, {false}}
	attackRadius          = []float64{1.33, 1.5, 1.8, 1.5, 0.8}
)

const normalHitNum = 5

func init() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 21)
	attackFrames[0][action.ActionAttack] = 15

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 24)
	attackFrames[1][action.ActionAttack] = 16

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 36)
	attackFrames[2][action.ActionAttack] = 27

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][1], 58)
	attackFrames[3][action.ActionAttack] = 31

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 66)
	attackFrames[4][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	centerTarget := c.Core.Combat.Player()
	if c.NormalCounter == 4 {
		centerTarget = c.Core.Combat.PrimaryTarget() // N5 is a bullet
	}
	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:          combat.AttackTagNormal,
			ICDTag:             combat.ICDTagNormalAttack,
			ICDGroup:           combat.ICDGroupDefault,
			StrikeType:         combat.StrikeTypeSlash,
			Element:            attributes.Physical,
			Durability:         25,
			Mult:               mult[c.TalentLvlAttack()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}
		radius := attackRadius[c.NormalCounter]
		c.QueueCharTask(func() {
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHit(centerTarget, radius),
				0,
				0,
			)
		}, attackHitmarks[c.NormalCounter][i])
	}

	if c.Base.Cons >= 6 {
		c.c6("attack")
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}
