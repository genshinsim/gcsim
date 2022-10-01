package candace

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	attackFrames   [][]int
	attackHitmarks = [][]int{{22}, {32}, {52, 52}, {94}} // TODO add correct hitmarks
)

const normalHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)                               // TODO: add correct frame data
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 22) // N1 -> CA
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 32) // N2 -> CA
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 52) // N3 -> CA
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][1], 94) // N4 -> CA
	attackFrames[3][action.ActionCharge] = 500                               // N5 -> CA, TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:          combat.AttackTagNormal,
			ICDTag:             combat.ICDTagNormalAttack,
			ICDGroup:           combat.ICDGroupPole,
			Element:            attributes.Physical,
			Durability:         25,
			Mult:               mult[c.TalentLvlAttack()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   0.03 * 60, // TODO: verify hitlag frames
			CanBeDefenseHalted: true,      // TODO: verify defense halt
		}
		c.QueueCharTask(func() {
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHit(c.Core.Combat.Player(), 0.1, false, combat.TargettableEnemy),
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
