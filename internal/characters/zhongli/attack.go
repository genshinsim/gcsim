package zhongli

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
	attackHits            = []int{1, 1, 1, 1, 4, 1}
	attackEarliestCancel  = []int{11, 9, 8, 16, 4, 29}
	attackHitmarks        = [][]int{{11}, {9}, {8}, {16}, {11, 18, 23, 29}, {29}}
	attackHitlagHaltFrame = [][]float64{{0.02}, {0.02}, {0.02}, {0.02}, {0, 0, 0, 0}, {0.02}}
	attackDefHalt         = [][]bool{{true}, {true}, {true}, {true}, {false, false, false, false}, {true}}
	attackHitboxes        = [][]float64{{1.5, 3.8}, {2}, {1, 1.5}, {1.7}, {1, 4}, {1, 4}}
	attackOffsets         = []float64{0, 0.8, 0.5, 1.8, -1, 0.2}
	attackFanAngles       = []float64{360, 180, 360, 360, 360, 360}
)

const normalHitNum = 6

func init() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackEarliestCancel[0], 30)
	attackFrames[0][action.ActionAttack] = 18

	attackFrames[1] = frames.InitNormalCancelSlice(attackEarliestCancel[1], 30)
	attackFrames[1][action.ActionAttack] = 13

	attackFrames[2] = frames.InitNormalCancelSlice(attackEarliestCancel[2], 28)
	attackFrames[2][action.ActionAttack] = 19

	attackFrames[3] = frames.InitNormalCancelSlice(attackEarliestCancel[3], 34)
	attackFrames[3][action.ActionCharge] = 33

	attackFrames[4] = frames.InitNormalCancelSlice(attackEarliestCancel[4], 31)
	attackFrames[4][action.ActionAttack] = 27
	attackFrames[4][action.ActionSkill] = 5
	attackFrames[4][action.ActionBurst] = 5
	attackFrames[4][action.ActionDash] = 5
	attackFrames[4][action.ActionJump] = 5

	attackFrames[5] = frames.InitNormalCancelSlice(attackEarliestCancel[5], 54)
	attackFrames[5][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	for i := 0; i < attackHits[c.NormalCounter]; i++ {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeSpear,
			Element:            attributes.Physical,
			Durability:         25,
			Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
			FlatDmg:            c.a4Attacks(),
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}
		if c.NormalCounter == 1 || c.NormalCounter == 4 {
			ai.StrikeType = attacks.StrikeTypeSlash
		}
		ap := combat.NewCircleHitOnTargetFanAngle(
			c.Core.Combat.Player(),
			geometry.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
			attackFanAngles[c.NormalCounter],
		)
		if c.NormalCounter == 0 || c.NormalCounter == 2 || c.NormalCounter == 4 || c.NormalCounter == 5 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: attackOffsets[c.NormalCounter]},
				attackHitboxes[c.NormalCounter][0],
				attackHitboxes[c.NormalCounter][1],
			)
		}
		// the multihit part generates no hitlag so this is fine
		c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter][i], attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackEarliestCancel[c.NormalCounter],
		State:           action.NormalAttackState,
	}, nil
}
