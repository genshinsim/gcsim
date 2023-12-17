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

// TODO frames & aoe
var (
	attackFrames   [][]int
	attackHitmarks = []int{0, 0, 0}
	attackRadius   = []float64{0, 0, 0}
	attackAngle    = []float64{0, 0, 0}
)

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 0)
	attackFrames[0][action.ActionAttack] = 0
	attackFrames[0][action.ActionCharge] = 0
	attackFrames[0][action.ActionSkill] = 0
	attackFrames[0][action.ActionBurst] = 0
	attackFrames[0][action.ActionDash] = 0
	attackFrames[0][action.ActionJump] = 0
	attackFrames[0][action.ActionSwap] = 0

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 0)
	attackFrames[1][action.ActionAttack] = 0
	attackFrames[1][action.ActionCharge] = 0
	attackFrames[1][action.ActionSkill] = 0
	attackFrames[1][action.ActionBurst] = 0
	attackFrames[1][action.ActionDash] = 0
	attackFrames[1][action.ActionJump] = 0
	attackFrames[1][action.ActionSwap] = 0

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 0)
	attackFrames[2][action.ActionAttack] = 0
	attackFrames[2][action.ActionCharge] = 0
	attackFrames[2][action.ActionSkill] = 0
	attackFrames[2][action.ActionBurst] = 0
	attackFrames[2][action.ActionDash] = 0
	attackFrames[2][action.ActionJump] = 0
	attackFrames[2][action.ActionSwap] = 0
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
