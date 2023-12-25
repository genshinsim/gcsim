package charlotte

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

const normalHitNum = 3

// aoe
var (
	attackFrames   [][]int
	attackHitmarks = []int{13, 25, 31}
	attackRadius   = []float64{0, 0, 0}
	attackAngle    = []float64{0, 0, 0}
)

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 31) // N1 -> Walk
	attackFrames[0][action.ActionAttack] = 23
	attackFrames[0][action.ActionCharge] = 24
	attackFrames[0][action.ActionSkill] = 6
	attackFrames[0][action.ActionBurst] = 5
	attackFrames[0][action.ActionDash] = 4
	attackFrames[0][action.ActionJump] = 5
	attackFrames[0][action.ActionSwap] = 13

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 44) // N2 -> Walk
	attackFrames[1][action.ActionAttack] = 35
	attackFrames[1][action.ActionCharge] = 25
	attackFrames[1][action.ActionSkill] = 14
	attackFrames[1][action.ActionBurst] = 16
	attackFrames[1][action.ActionDash] = 15
	attackFrames[1][action.ActionJump] = 15
	attackFrames[1][action.ActionSwap] = 17

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 69) // N3 -> Walk
	attackFrames[2][action.ActionAttack] = 74
	attackFrames[2][action.ActionCharge] = 65
	attackFrames[2][action.ActionSkill] = 13
	attackFrames[2][action.ActionBurst] = 6
	attackFrames[2][action.ActionDash] = 7
	attackFrames[2][action.ActionJump] = 6
	attackFrames[2][action.ActionSwap] = 5
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     attacks.ICDTagNormalAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	ap := combat.NewCircleHitOnTargetFanAngle(
		c.Core.Combat.Player(),
		nil,
		attackRadius[c.NormalCounter],
		attackAngle[c.NormalCounter],
	)

	c.Core.QueueAttack(
		ai,
		ap,
		attackHitmarks[c.NormalCounter],
		attackHitmarks[c.NormalCounter],
	)

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}, nil
}
