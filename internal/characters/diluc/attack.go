package diluc

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	attackFrames          [][]int
	attackHitmarks        = []int{24, 39, 26, 49}
	attackHitlagHaltFrame = []float64{.1, .09, .09, .12}
	attackRadius          = []float64{2, 1.8, 2, 1.8}
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 32)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 46)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 34)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 99)
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	// C6: Additionally, Searing Onslaught will not interrupt the Normal Attack combo.
	if c.Base.Cons >= 6 && c.Core.Player.CurrentState() == action.SkillState {
		c.NormalCounter = c.savedNormalCounter
	}

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:          combat.AttackTagNormal,
		ICDTag:             combat.ICDTagNormalAttack,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeBlunt,
		Element:            attributes.Physical,
		Durability:         25,
		Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter] * 60,
		CanBeDefenseHalted: true,
	}
	radius := attackRadius[c.NormalCounter]
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), radius),
		attackHitmarks[c.NormalCounter],
		attackHitmarks[c.NormalCounter],
	)

	defer func() {
		c.AdvanceNormalIndex()
		c.savedNormalCounter = c.NormalCounter
	}()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}
