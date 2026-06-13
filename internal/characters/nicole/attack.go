package nicole

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	attackFrames   [][]int
	attackHitmarks = []int{15, 8, 38}
)

const (
	normalHitNum = 3
)

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 25) // N1 -> N2
	attackFrames[0][action.ActionCharge] = 18                             // N1 -> CA

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 34) // N2 -> W
	attackFrames[1][action.ActionAttack] = 22                             // N2 -> N3
	attackFrames[1][action.ActionCharge] = 17                             // N2 -> CA

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 67) // N3 -> W
	attackFrames[2][action.ActionAttack] = 52                             // N3 -> N1
	attackFrames[2][action.ActionCharge] = 52                             // N3 -> CA
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	// TODO: Find actual windup and conditions for windup
	windup := 0
	switch c.Core.Player.CurrentState() {
	case action.NormalAttackState:
		if c.NormalCounter == 0 {
			// N3 -> N1
			windup = 10
		}
	case action.ChargeAttackState:
		windup = -4
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     attacks.ICDTagNormalAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			nil,
			0.75,
		),
		attackHitmarks[c.NormalCounter]+windup,
		attackHitmarks[c.NormalCounter]+windup,
	)

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          attackFuncWithWindup(c.Character, attackFrames, windup),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}, nil
}

func attackFuncWithWindup(c *character.Character, slice [][]int, windup int) func(action.Action) int {
	n := c.NormalCounter
	atkspd := c.Stat(attributes.AtkSpd)
	return func(next action.Action) int {
		return frames.AtkSpdAdjust(slice[n][next], atkspd) + windup
	}
}
