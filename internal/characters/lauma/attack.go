package lauma

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
	attackFrames [][]int

	attackHitmarks = []int{14, 11, 16}
	attackOffsets  = []float64{0, 0, 0}
	attackHitboxes = [][]float64{{2.5, 8}, {2.5, 8}, {4, 8}}
)

const normalHitNum = 3

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 26) // N1 -> Walk
	attackFrames[0][action.ActionAttack] = 29
	attackFrames[0][action.ActionCharge] = 17

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 33) // N2 -> N3
	attackFrames[1][action.ActionCharge] = 19
	attackFrames[1][action.ActionWalk] = 28

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 38) // N3 -> N1
	attackFrames[2][action.ActionCharge] = 33
	attackFrames[2][action.ActionWalk] = 38
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     attacks.ICDTagNormalAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			info.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
			attackHitboxes[c.NormalCounter][1],
		),
		attackHitmarks[c.NormalCounter],
		attackHitmarks[c.NormalCounter],
	)

	if c.Base.Cons >= 6 {
		aiC6 := info.AttackInfo{
			ActorIndex:       c.Index(),
			Abil:             "Normal C6 Pale Hymn",
			AttackTag:        attacks.AttackTagDirectLunarBloom,
			ICDTag:           attacks.ICDTagNormalAttack,
			ICDGroup:         attacks.ICDGroupDefault,
			StrikeType:       attacks.StrikeTypeDefault,
			Element:          attributes.Dendro,
			Durability:       25,
			UseEM:            true,
			Mult:             1.5,
			IgnoreDefPercent: 1,
		}

		paleHymnStacks--

		c.Core.QueueAttack(
			aiC6,
			combat.NewBoxHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				info.Point{Y: attackOffsets[c.NormalCounter]},
				attackHitboxes[c.NormalCounter][0],
				attackHitboxes[c.NormalCounter][1],
			),
			attackHitmarks[c.NormalCounter],
			attackHitmarks[c.NormalCounter],
		)
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}, nil
}
