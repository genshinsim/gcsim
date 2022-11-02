package heizou

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = [][]int{{12}, {13}, {21}, {13, 19, 27}, {31}}
var attackHitlagHaltFrame = [][]float64{{0.03}, {0.03}, {0.06}, {0, 0, 0.09}, {0.12}}
var attackDefHalt = [][]bool{{true}, {true}, {true}, {false, false, true}, {true}}

const normalHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 21)
	attackFrames[0][action.ActionAttack] = 20
	attackFrames[0][action.ActionCharge] = 21

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 21)
	attackFrames[1][action.ActionAttack] = 17
	attackFrames[1][action.ActionCharge] = 21

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 46)
	attackFrames[2][action.ActionAttack] = 45
	attackFrames[2][action.ActionCharge] = 46

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][2], 38)
	attackFrames[3][action.ActionAttack] = 36
	attackFrames[3][action.ActionCharge] = 38

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 66)
	attackFrames[4][action.ActionAttack] = 66
	attackFrames[4][action.ActionCharge] = 500
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	for i := 0; i < hits[c.NormalCounter]; i++ {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:          combat.AttackTagNormal,
			ICDTag:             combat.ICDTagNormalAttack,
			ICDGroup:           combat.ICDGroupDefault,
			Element:            attributes.Anemo,
			Durability:         25,
			Mult:               attack[c.NormalCounter][i][c.TalentLvlAttack()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}
		// multihit on N4 only has hitlag on last hit so no need for char queue here
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), 0.1),
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
