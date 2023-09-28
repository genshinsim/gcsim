package kokomi

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = []int{4, 12, 28}

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
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	ai.FlatDmg = c.burstDmgBonus(ai.AttackTag)

	radius := 0.7
	if c.Core.Status.Duration(burstKey) > 0 {
		radius = 1.2
	}

	// TODO: Assume that this is not dynamic (snapshot on projectile release)
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, radius),
		attackHitmarks[c.NormalCounter],
		attackHitmarks[c.NormalCounter]+travel,
		c.makeBurstHealCB(),
		c.makeC4CB(),
	)
	if c.NormalCounter == c.NormalHitNum-1 {
		c.c1(attackHitmarks[c.NormalCounter], travel)
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}, nil
}
