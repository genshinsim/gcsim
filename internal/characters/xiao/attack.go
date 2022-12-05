package xiao

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	attackFrames          [][]int
	attackHitmarks        = [][]int{{4, 17}, {15}, {15}, {14, 31}, {16}, {39}}
	attackHitlagHaltFrame = [][]float64{{0, 0.01}, {0.01}, {0.01}, {0.02, 0.02}, {0.02}, {0.04}}
	attackDefHalt         = [][]bool{{false, true}, {true}, {true}, {false, true}, {true}, {true}}
	attackRadius          = [][][]float64{{{1.8, 1.52}, {1.6}, {1.6}, {1.6, 1.8}, {1.68}, {2}}, {{2, 1.7}, {1.8}, {1.8}, {1.8, 2}, {1.81}, {2.4}}}
	attackStrikeTypes     = [][]combat.StrikeType{
		{combat.StrikeTypeSlash, combat.StrikeTypeSpear},
		{combat.StrikeTypeSlash},
		{combat.StrikeTypeSlash},
		{combat.StrikeTypeSlash, combat.StrikeTypeSlash},
		{combat.StrikeTypeSpear},
		{combat.StrikeTypeSlash},
	}
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
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}

		burstIndex := 0
		if c.StatusIsActive(burstBuffKey) {
			burstIndex = 1
		}
		radius := attackRadius[burstIndex][c.NormalCounter][i]
		c.QueueCharTask(func() {
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHit(c.Core.Combat.Player(), radius),
				0,
				0,
			)
		}, attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}
