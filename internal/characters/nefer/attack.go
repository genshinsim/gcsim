package nefer

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

const normalHitNum = 4

func normalAttackCanQueueAfter(normalIndex int) int {
	return attackFrames[normalIndex][action.ActionSkill]
}

var (
	attackFrames           [][]int
	attackHitmarks         = [][]int{{10}, {8}, {10, 10 + 16}, {22}}
	attackHitlagHaltFrames = []float64{0.035, 0.035, 0.03, 0.05}
	attackHitboxes         = [][][]float64{{{1.5}}, {{2, 3.5}}, {{2.1}, {2.5, 3}}, {{2.5}}}
	attackOffsets          = [][]float64{{0.5}, {0}, {1.2, -0.3}, {1}}
)

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 40)
	attackFrames[0][action.ActionAttack] = 15
	attackFrames[0][action.ActionWalk] = 33

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 38)
	attackFrames[1][action.ActionAttack] = 20
	attackFrames[1][action.ActionWalk] = 34

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 62)
	attackFrames[2][action.ActionAttack] = 49
	attackFrames[2][action.ActionWalk] = 50

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 63)
	attackFrames[3][action.ActionAttack] = 51
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex:       c.Index(),
		Abil:             fmt.Sprintf("Normal %v", c.NormalCounter+1),
		AttackTag:        attacks.AttackTagNormal,
		ICDTag:           attacks.ICDTagNormalAttack,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Dendro,
		Durability:       25,
		HitlagFactor:     0.01,
		HitlagHaltFrames: attackHitlagHaltFrames[c.NormalCounter],
		Mult:             attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	var cb info.AttackPattern

	if c.NormalCounter == 1 {
		cb = combat.NewBoxHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			info.Point{Y: attackOffsets[c.NormalCounter][0]},
			attackHitboxes[c.NormalCounter][0][0],
			attackHitboxes[c.NormalCounter][0][1],
		)
	} else {
		cb = combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			info.Point{Y: attackOffsets[c.NormalCounter][0]},
			attackHitboxes[c.NormalCounter][0][0],
		)
	}

	if c.NormalCounter == 2 {
		c.Core.QueueAttack(ai, cb, attackHitmarks[c.NormalCounter][0], attackHitmarks[c.NormalCounter][0])

		cb = combat.NewBoxHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			info.Point{Y: attackOffsets[c.NormalCounter][0]},
			attackHitboxes[c.NormalCounter][1][0],
			attackHitboxes[c.NormalCounter][1][1],
		)
		c.Core.QueueAttack(ai, cb, attackHitmarks[c.NormalCounter][1], attackHitmarks[c.NormalCounter][1])
	} else {
		c.Core.QueueAttack(ai, cb, attackHitmarks[c.NormalCounter][0], attackHitmarks[c.NormalCounter][0])
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   normalAttackCanQueueAfter(c.NormalCounter),
		State:           action.NormalAttackState,
	}, nil
}
