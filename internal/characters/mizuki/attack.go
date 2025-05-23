package mizuki

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = []int{7, 19, 37}

const normalHitNum = 3

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 10) // Common frames
	attackFrames[0][action.ActionAttack] = 18                             // N1 -> N2
	attackFrames[0][action.ActionCharge] = 20
	attackFrames[0][action.ActionWalk] = 34
	attackFrames[0][action.ActionSwap] = 13

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 21) // Common frames
	attackFrames[1][action.ActionAttack] = 37                             // N2 -> N3
	attackFrames[1][action.ActionCharge] = 36
	attackFrames[1][action.ActionWalk] = 38

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 41) // Common frames
	attackFrames[2][action.ActionCharge] = 98                             // N3 -> N1
	attackFrames[2][action.ActionWalk] = 72
}

// Standard attack damage function
// Has "travel" parameter, used to set the number of frames that the projectile is in the air (default = 10)
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
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	radius := 0.5

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, radius),
		attackHitmarks[c.NormalCounter],
		attackHitmarks[c.NormalCounter]+travel,
	)

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}, nil
}
