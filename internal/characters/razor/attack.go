package razor

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var (
	attackFrames          [][]int
	attackHitmarks        = []int{25, 16, 13, 38}
	attackPoiseDMG        = []float64{112.93, 97.29, 121.67, 160.19}
	attackHitlagHaltFrame = []float64{0.1, 0.1, 0.1, 0.15}
	attackHitlagFactor    = []float64{0.01, 0.01, 0.05, 0.01}
	attackHitboxes        = [][]float64{{2}, {3.2, 3}, {2}, {2}}
	attackOffsets         = []float64{1, 0.5, 1, 1.8}
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum) // should be 4

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 45)  // N1 -> N2
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 33)  // N2 -> N3
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 47)  // N3 -> N4
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 116) // N4 -> N1
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
		AttackTag:          attacks.AttackTagNormal,
		ICDTag:             attacks.ICDTagNormalAttack,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		PoiseDMG:           attackPoiseDMG[c.NormalCounter],
		Element:            attributes.Physical,
		Durability:         25,
		HitlagFactor:       attackHitlagFactor[c.NormalCounter],
		HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter],
		CanBeDefenseHalted: true,
	}
	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		geometry.Point{Y: attackOffsets[c.NormalCounter]},
		attackHitboxes[c.NormalCounter][0],
	)
	if c.NormalCounter == 1 {
		ap = combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
			attackHitboxes[c.NormalCounter][1],
		)
	}

	var cb combat.AttackCBFunc
	if c.StatusIsActive(burstBuffKey) {
		cb = c.wolfBurst(c.NormalCounter)
	}
	var c6cb func(a combat.AttackCB)
	if c.Base.Cons >= 6 {
		c6cb = c.c6cb
	}
	c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter], attackHitmarks[c.NormalCounter], cb, c6cb)

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}, nil
}
