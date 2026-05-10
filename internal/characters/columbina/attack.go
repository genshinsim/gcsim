package columbina

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
	attackHitmarks = []int{23, 20, 41}
	attackRadius   = []float64{1.0, 1.0, 2.5}
)

const normalHitNum = 3

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 36)
	attackFrames[0][action.ActionAttack] = 24
	attackFrames[0][action.ActionCharge] = 24

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 37)
	attackFrames[1][action.ActionAttack] = 20
	attackFrames[1][action.ActionCharge] = 33

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 76)
	attackFrames[2][action.ActionAttack] = 74
	attackFrames[2][action.ActionCharge] = 71
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	windup := 0
	if c.NormalCounter == 0 {
		switch c.Core.Player.CurrentState() {
		case action.NormalAttackState:
			windup = 2
		case action.ChargeAttackState:
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

	radius := attackRadius[c.NormalCounter]

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, radius),
		attackHitmarks[c.NormalCounter]+windup,
		attackHitmarks[c.NormalCounter]+windup,
	)
	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          c.newAttackFunc(attackFrames, windup),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}, nil
}

func (c *char) newAttackFunc(slice [][]int, windup int) func(action.Action) int {
	n := c.NormalCounter
	atkspd := c.Stat(attributes.AtkSpd)
	return func(next action.Action) int {
		return frames.AtkSpdAdjust(slice[n][next], atkspd) + windup
	}
}
