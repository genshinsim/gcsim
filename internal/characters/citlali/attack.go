package citlali

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

const normalHitNum = 3

var (
	attackFrames         [][]int
	attackHitmarks       = []int{21, 25, 27}
	attackEarliestCancel = []int{5, 15, 4}
	attackRadius         = []float64{1, 1, 2}
)

// charlotte frames. CHANGE
func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 32) // N1 -> Walk
	attackFrames[0][action.ActionAttack] = 23
	attackFrames[0][action.ActionCharge] = 25
	attackFrames[0][action.ActionSkill] = 5
	attackFrames[0][action.ActionBurst] = 5
	attackFrames[0][action.ActionDash] = 5
	attackFrames[0][action.ActionJump] = 6
	attackFrames[0][action.ActionSwap] = 23

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 45) // N2 -> Walk
	attackFrames[1][action.ActionAttack] = 35
	attackFrames[1][action.ActionCharge] = 26
	attackFrames[1][action.ActionSkill] = 15
	attackFrames[1][action.ActionBurst] = 17
	attackFrames[1][action.ActionDash] = 16
	attackFrames[1][action.ActionJump] = 16
	attackFrames[1][action.ActionSwap] = 16

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 74) // N3 -> N1
	attackFrames[2][action.ActionCharge] = 66
	attackFrames[2][action.ActionSkill] = 14
	attackFrames[2][action.ActionBurst] = 7
	attackFrames[2][action.ActionDash] = 8
	attackFrames[2][action.ActionJump] = 7
	attackFrames[2][action.ActionWalk] = 70
	attackFrames[2][action.ActionSwap] = 4
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 1
	}

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

	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.PrimaryTarget(),
		nil,
		attackRadius[c.NormalCounter],
	)

	c.Core.QueueAttack(
		ai,
		ap,
		attackHitmarks[c.NormalCounter]+travel,
		attackHitmarks[c.NormalCounter]+travel,
	)

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackEarliestCancel[c.NormalCounter] + travel,
		State:           action.NormalAttackState,
	}, nil
}
