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
	attackHitmarks = []int{4, 12, 28}
	attackRadius   = []float64{1.0, 1.0, 2.5}
)

const normalHitNum = 3

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 30)
	attackFrames[0][action.ActionAttack] = 14
	attackFrames[0][action.ActionCharge] = 19

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 34)
	attackFrames[1][action.ActionAttack] = 30

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 65)
	attackFrames[2][action.ActionCharge] = 60
	attackFrames[2][action.ActionWalk] = 60
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
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
		attackHitmarks[c.NormalCounter],
		attackHitmarks[c.NormalCounter],
	)
	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}, nil
}
