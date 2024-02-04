package navia

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
	attackHitmarks        = [][]int{{23}, {22}, {31, 39, 48}, {41}}
	attackPoiseDMG        = []float64{129.8, 120.1, 50.5, 185.9}
	attackHitlagHaltFrame = [][]float64{{0.06}, {0.06}, {0.01, 0.01, 0.01}, {0.06}}
	attackDefHalt         = [][]bool{{true}, {true}, {false, false, false}, {true}}
	attackHitboxes        = [][]float64{{2}, {2, 4.3}, {3, 4.5}, {2, 4.7}}
	attackOffsets         = []float64{0.5, -1.5, 0.3, -1.85}
	attackEarliestCancel  = []int{23, 22, 29, 41}
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackEarliestCancel[0], 28)
	attackFrames[1] = frames.InitNormalCancelSlice(attackEarliestCancel[1], 42)
	attackFrames[2] = frames.InitNormalCancelSlice(attackEarliestCancel[2], 48)
	attackFrames[2][action.ActionSkill] = 30
	attackFrames[3] = frames.InitNormalCancelSlice(attackEarliestCancel[3], 93)
}

func (c *char) Attack(_ map[string]int) (action.Info, error) {
	for i, mult := range auto[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeBlunt,
			PoiseDMG:           attackPoiseDMG[c.NormalCounter],
			Element:            attributes.Physical,
			Durability:         25,
			Mult:               mult[c.TalentLvlAttack()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
			IsDeployable:       c.NormalCounter == 2,
		}
		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
		)
		if c.NormalCounter != 0 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: attackOffsets[c.NormalCounter]},
				attackHitboxes[c.NormalCounter][0],
				attackHitboxes[c.NormalCounter][1],
			)
		}
		// no char queue is fine here because multhit doesn't have hitlag on her
		// N3 should snap on gadget creation which is assumed to be earliest cancel here so that infusion is snapped properly
		c.Core.QueueAttack(ai, ap, attackEarliestCancel[c.NormalCounter], attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackEarliestCancel[c.NormalCounter],
		State:           action.NormalAttackState,
	}, nil
}
