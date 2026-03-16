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
	attackFrames   [][]int
	attackHitmarks = []int{10, 8, 26, 22}
	attackHitboxes = [][]float64{{2, 8}, {2, 8}, {2.5, 9}, {2.8, 8}}
	attackOffsets  = []float64{0, 0, 0, -0.5}
)

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 40)
	attackFrames[0][action.ActionAttack] = 15
	attackFrames[0][action.ActionCharge] = 40
	attackFrames[0][action.ActionWalk] = 33

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 38)
	attackFrames[1][action.ActionAttack] = 20
	attackFrames[1][action.ActionCharge] = 38
	attackFrames[1][action.ActionWalk] = 34

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 62)
	attackFrames[2][action.ActionAttack] = 49
	attackFrames[2][action.ActionCharge] = 62
	attackFrames[2][action.ActionWalk] = 50

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 51)
	attackFrames[3][action.ActionWalk] = 63
	attackFrames[3][action.ActionAttack] = 51
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter+1),
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     attacks.ICDTagNormalAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	cb := combat.NewBoxHit(
		c.Core.Combat.Player(),
		c.Core.Combat.PrimaryTarget(),
		info.Point{Y: attackOffsets[c.NormalCounter]},
		attackHitboxes[c.NormalCounter][0],
		attackHitboxes[c.NormalCounter][1],
	)

	if c.NormalCounter == 2 {
		c.Core.QueueAttack(ai, cb, 10, 10)
		c.Core.QueueAttack(ai, cb, 26, 26)
	} else {
		c.Core.QueueAttack(ai, cb, attackHitmarks[c.NormalCounter], attackHitmarks[c.NormalCounter])
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   normalAttackCanQueueAfter(c.NormalCounter),
		State:           action.NormalAttackState,
	}, nil
}

func (c *char) HighPlungeAttack(map[string]int) (action.Info, error) {
	return c.plungeAttack(high[0][c.TalentLvlAttack()])
}

func (c *char) LowPlungeAttack(map[string]int) (action.Info, error) {
	return c.plungeAttack(low[0][c.TalentLvlAttack()])
}

func (c *char) plungeAttack(mult float64) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Plunge",
		AttackTag:  attacks.AttackTagPlunge,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       mult,
	}

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 4.5), 1, 1)

	return action.Info{State: action.PlungeAttackState}, nil
}
