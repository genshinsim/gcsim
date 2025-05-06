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
	attackFrames   [][]int
	attackHitmarks = []int{16, 16, 36}
	attackRadius   = []float64{0.75, 0.75, 0.75}
)

// charlotte frames. CHANGE
func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], attackHitmarks[0]) // N1 -> Earliest cancel (jump)
	attackFrames[0][action.ActionAttack] = 34
	attackFrames[0][action.ActionCharge] = 33
	attackFrames[0][action.ActionWalk] = 38

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], attackHitmarks[1]) // N2 -> Earliest cancel (jump)
	attackFrames[1][action.ActionAttack] = 37
	attackFrames[1][action.ActionCharge] = 37
	attackFrames[1][action.ActionWalk] = 39

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], attackHitmarks[2]) // N3 -> Earliest cancel (jump)
	attackFrames[0][action.ActionAttack] = 49
	attackFrames[2][action.ActionCharge] = 50
	attackFrames[1][action.ActionWalk] = 52
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
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
		CanQueueAfter:   attackHitmarks[c.NormalCounter] + travel,
		State:           action.NormalAttackState,
	}, nil
}
