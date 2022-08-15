package barbara

import (
	"fmt"
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var (
	attackFrames          [][]int
	attackFramesWithLag   [][]int
	attackHitmarks        = []int{6, 11, 12, 32}
	attackHitmarksWithLag []int
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	// N1 -> x
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 23)
	attackFrames[0][action.ActionAttack] = 15
	attackFrames[0][action.ActionCharge] = 18

	// N2 -> x
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 28)
	attackFrames[1][action.ActionAttack] = 21
	attackFrames[1][action.ActionCharge] = 24

	// N3 -> x
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 30)
	attackFrames[2][action.ActionAttack] = 22
	attackFrames[2][action.ActionCharge] = 28

	// N4 -> x
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 60)
	attackFrames[3][action.ActionWalk] = 57
	attackFrames[3][action.ActionDash] = 2
	attackFrames[3][action.ActionJump] = 3
	attackFrames[3][action.ActionSwap] = 2

	// N1 -> x (Dash/N4 -> N1 8f lag)
	attackFramesWithLag = make([][]int, len(attackFrames))
	for i := range attackFrames {
		attackFramesWithLag[i] = make([]int, len(attackFrames[i]))
		copy(attackFramesWithLag[i], attackFrames[i])
	}
	for i := range attackFramesWithLag[0] {
		attackFramesWithLag[0][i] += 8
	}

	// Hitmarks (Dash/N4 -> N1 8f lag)
	attackHitmarksWithLag = make([]int, len(attackHitmarks))
	for i := range attackHitmarks {
		attackHitmarksWithLag[i] = attackHitmarks[i] + 8
	}
}

// Standard attack function with seal handling
func (c *char) Attack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  combat.AttackTagNormal,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	done := false
	cb := func(_ combat.AttackCB) {
		if done {
			return
		}
		// check for healing
		if c.Core.Status.Duration(barbSkillKey) > 0 {
			c.Core.Player.Heal(player.HealInfo{
				Caller:  c.Index,
				Target:  -1,
				Message: "Melody Loop (Normal Attack)",
				Src:     prochpp[c.TalentLvlSkill()]*c.MaxHP() + prochp[c.TalentLvlSkill()],
				Bonus:   c.Stat(attributes.Heal),
			})
			done = true
		}
	}

	adjustedFrames := attackFrames
	adjustedHitmarks := attackHitmarks
	currState := c.Core.Player.CurrentState()
	if currState == action.DashState || (currState == action.NormalAttackState && c.NormalCounter == 0) {
		adjustedFrames = attackFramesWithLag
		adjustedHitmarks = attackHitmarksWithLag
	}
	c.Core.QueueAttack(ai,
		combat.NewCircleHit(c.Core.Combat.Player(), 0.1, false, combat.TargettableEnemy),
		adjustedHitmarks[c.NormalCounter],
		adjustedHitmarks[c.NormalCounter],
		cb,
	)

	defer c.AdvanceNormalIndex()

	canQueueAfter := math.MaxInt32
	for _, f := range adjustedFrames[c.NormalCounter] {
		if f < canQueueAfter {
			canQueueAfter = f
		}
	}
	// return animation cd
	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, adjustedFrames),
		AnimationLength: adjustedFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   canQueueAfter,
		State:           action.NormalAttackState,
	}
}
