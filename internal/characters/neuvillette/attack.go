package neuvillette

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

const normalHitNum = 3

var (
	attackFrames   [][]int
	attackHitmarks = []int{19, 16, 32}
	attackHitboxes = []float64{1.0, 1.0, 1.5}
)

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 36)
	attackFrames[0][action.ActionAttack] = 29
	attackFrames[0][action.ActionCharge] = 20

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 33)
	attackFrames[1][action.ActionAttack] = 31
	attackFrames[1][action.ActionCharge] = 22

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 62)
	attackFrames[2][action.ActionWalk] = 61
	attackFrames[2][action.ActionCharge] = 51
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.chargeEarlyCancelled {
		return action.Info{}, fmt.Errorf("%v: Cannot early cancel Charged Attack: Equitable Judgement with Normal Attack", c.CharWrapper.Base.Key)
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     attacks.ICDTagNormalAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(
			c.Core.Combat.PrimaryTarget(),
			geometry.Point{},
			attackHitboxes[c.NormalCounter],
		),
		attackHitmarks[c.NormalCounter],
		attackHitmarks[c.NormalCounter],
	)

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackFrames[c.NormalCounter][action.ActionSwap],
		State:           action.NormalAttackState,
	}, nil
}
