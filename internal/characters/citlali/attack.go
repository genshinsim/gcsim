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
	attackFrames       [][]int
	attackHitmarks     = []int{16, 16, 36}
	attackRadius       = 0.8
	attackHitlagFactor = 0.05
)

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 38) // N1 -> Walk
	attackFrames[0][action.ActionAttack] = 34
	attackFrames[0][action.ActionCharge] = 33

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 39) // N2 -> Walk
	attackFrames[1][action.ActionAttack] = 37
	attackFrames[1][action.ActionCharge] = 37

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 52) // N3 -> Walk
	attackFrames[2][action.ActionAttack] = 49
	attackFrames[2][action.ActionCharge] = 50
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	ai := combat.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:    attacks.AttackTagNormal,
		ICDTag:       attacks.ICDTagNormalAttack,
		ICDGroup:     attacks.ICDGroupDefault,
		StrikeType:   attacks.StrikeTypeDefault,
		Element:      attributes.Cryo,
		Durability:   25,
		Mult:         attack[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor: attackHitlagFactor,
	}

	ap := combat.NewCircleHit(
		c.Core.Combat.Player(),
		c.Core.Combat.PrimaryTarget(),
		nil,
		attackRadius,
	)

	c.Core.QueueAttack(
		ai,
		ap,
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
