package mizuki

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
	attackHitmarks = []int{7, 19, 37}
)

const normalHitNum = 3

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 34) // N1 -> Walk
	attackFrames[0][action.ActionAttack] = 18
	attackFrames[0][action.ActionCharge] = 20
	attackFrames[0][action.ActionSkill] = 10
	attackFrames[0][action.ActionBurst] = 9
	attackFrames[0][action.ActionDash] = 10
	attackFrames[0][action.ActionSwap] = 13

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 38) // N2 -> Walk
	attackFrames[1][action.ActionAttack] = 37
	attackFrames[1][action.ActionCharge] = 36
	attackFrames[1][action.ActionSkill] = 21
	attackFrames[1][action.ActionBurst] = 20
	attackFrames[1][action.ActionDash] = 21
	attackFrames[1][action.ActionJump] = 21
	attackFrames[1][action.ActionSwap] = 21

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 98) // N3 -> CA
	attackFrames[2][action.ActionAttack] = 73
	attackFrames[2][action.ActionSkill] = 41
	attackFrames[2][action.ActionBurst] = 40
	attackFrames[2][action.ActionDash] = 41
	attackFrames[2][action.ActionJump] = 41
	attackFrames[2][action.ActionWalk] = 72
	attackFrames[2][action.ActionSwap] = 41
}

// Standard attack damage function
// Has "travel" parameter, used to set the number of frames that the projectile is in the air (default = 10)
func (c *char) Attack(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	ai := info.AttackInfo{
		ActorIndex:   c.Index(),
		Abil:         fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:    attacks.AttackTagNormal,
		ICDTag:       attacks.ICDTagNormalAttack,
		ICDGroup:     attacks.ICDGroupDefault,
		StrikeType:   attacks.StrikeTypeDefault,
		Element:      attributes.Anemo,
		Durability:   25,
		Mult:         attack[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor: 0.05,
		IsDeployable: true,
	}

	radius := 0.7

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
