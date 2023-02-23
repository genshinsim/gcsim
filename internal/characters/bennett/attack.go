package bennett

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	attackFrames           [][]int
	attackHitmarks         = []int{13, 9, 13, 25, 24}
	attackHitlagHaltFrames = []float64{0.03, 0.03, 0.06, 0.09, 0.12}
	attackHitboxes         = [][]float64{{1.2}, {1.2}, {2}, {1, 3.5}, {2}}
	attackOffsets          = []float64{0.8, 0.8, 0.6, 0.3, 0.8}
	attackFanAngles        = []float64{360, 360, 30, 360, 360}
)

const normalHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 33)
	attackFrames[0][action.ActionAttack] = 20

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 27)
	attackFrames[1][action.ActionAttack] = 17

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 46)
	attackFrames[2][action.ActionAttack] = 37

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 48)
	attackFrames[3][action.ActionAttack] = 44

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4], 60)
	attackFrames[4][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:          attacks.AttackTagNormal,
		ICDTag:             attacks.ICDTagNormalAttack,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeSlash,
		Element:            attributes.Physical,
		Durability:         25,
		HitlagHaltFrames:   attackHitlagHaltFrames[c.NormalCounter] * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
		Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	ap := combat.NewCircleHitOnTargetFanAngle(
		c.Core.Combat.Player(),
		combat.Point{Y: attackOffsets[c.NormalCounter]},
		attackHitboxes[c.NormalCounter][0],
		attackFanAngles[c.NormalCounter],
	)
	if c.NormalCounter == 3 {
		ap = combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			combat.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
			attackHitboxes[c.NormalCounter][1],
		)
	}
	c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter], attackHitmarks[c.NormalCounter])

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}
