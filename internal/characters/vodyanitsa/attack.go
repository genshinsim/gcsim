package vodyanitsa

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	attackFrames   [][]int
	attackHitmarks = []int{13, 13, 23, 34} // N1 hitmark includes windup
	attackSize     = [][]float64{{0.7}, {0.7}, {0.7}, {2, 10}}
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	// N1 -> x
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 26)
	attackFrames[0][action.ActionAttack] = 20
	attackFrames[0][action.ActionWalk] = 24

	// N2 -> x
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 34)
	attackFrames[1][action.ActionAttack] = 30
	attackFrames[1][action.ActionCharge] = 28

	// N3 -> x
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 48)
	attackFrames[2][action.ActionAttack] = 37
	attackFrames[2][action.ActionCharge] = 36

	// N4 -> x
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 105)
	attackFrames[3][action.ActionAttack] = 67
	attackFrames[3][action.ActionWalk] = 99
}

// Standard attack function with seal handling
func (c *char) Attack(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	windup := 0
	if c.NormalCounter == 0 {
		switch c.Core.Player.CurrentState() {
		case action.ChargeAttackState, action.NormalAttackState, action.SkillState:
			windup = 9
		case action.BurstState:
			windup = 2
		}
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     attacks.ICDTagNormalAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	ap := combat.NewCircleHit(
		c.Core.Combat.Player(),
		c.Core.Combat.PrimaryTarget(),
		nil,
		attackSize[c.NormalCounter][0],
	)

	if c.NormalCounter == 3 {
		ap = combat.NewBoxHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			nil,
			attackSize[c.NormalCounter][0],
			attackSize[c.NormalCounter][1],
		)
	}

	c.Core.QueueAttack(
		ai,
		ap,
		attackHitmarks[c.NormalCounter]+windup,
		attackHitmarks[c.NormalCounter]+windup+travel,
	)

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFuncWithOffset(c.Character, attackFrames, windup),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction] + windup,
		CanQueueAfter:   attackHitmarks[c.NormalCounter] + windup,
		State:           action.NormalAttackState,
	}, nil
}
