package sucrose

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
	attackHitmarks = []int{17, 18, 28, 28}
	attackRadius   = []float64{1, 1, 1, 2}
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	// TODO: check if hitmarks for NA->CA and CA->CA lines up
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 28) // walk
	attackFrames[0][action.ActionAttack] = 17
	attackFrames[0][action.ActionCharge] = 20

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 38) // walk
	attackFrames[1][action.ActionAttack] = 26
	attackFrames[1][action.ActionCharge] = 18

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 46) // walk
	attackFrames[2][action.ActionAttack] = 33
	attackFrames[2][action.ActionCharge] = 28

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[2], 65) // walk
	attackFrames[3][action.ActionAttack] = 51
	attackFrames[3][action.ActionCharge] = 54
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     attacks.ICDTagNormalAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	var c4cb info.AttackCBFunc
	if c.Base.Cons >= 4 {
		c4cb = c.makeC4Callback()
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			nil,
			attackRadius[c.NormalCounter],
		),
		attackHitmarks[c.NormalCounter],
		attackHitmarks[c.NormalCounter],
		c4cb,
	)

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}, nil
}
