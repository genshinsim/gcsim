package chongyun

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	attackFrames          [][]int
	attackHitmarks        = []int{26, 24, 41, 53}
	attackHitlagHaltFrame = []float64{.1, .09, .12, .12}
	attackHitboxes        = [][]float64{{2}, {2}, {2}, {2, 3}}
	attackOffsets         = []float64{1, 1, 1, -0.5}
	attackFanAngles       = []float64{360, 270, 360, 360}
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 30)  // N1 -> N2
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 36)  // N2 -> N3
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 57)  // N3 -> N4
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 101) // N4 -> N1
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		ActorIndex:         c.Index,
		AttackTag:          attacks.AttackTagNormal,
		ICDTag:             attacks.ICDTagNormalAttack,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		Element:            attributes.Physical,
		Durability:         25,
		Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter] * 60,
		CanBeDefenseHalted: true,
	}
	ap := combat.NewCircleHitOnTargetFanAngle(
		c.Core.Combat.Player(),
		combat.Point{Y: attackOffsets[c.NormalCounter]},
		attackHitboxes[c.NormalCounter][0],
		attackFanAngles[c.NormalCounter],
	)
	if c.NormalCounter == 3 {
		ap = combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			combat.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
			attackHitboxes[c.NormalCounter][1],
		)
	}
	c4CB := c.makeC4Callback()
	c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter], attackHitmarks[c.NormalCounter], c4CB)

	if c.Base.Cons >= 1 && c.NormalCounter == 3 {
		ai := combat.AttackInfo{
			Abil:       "Chongyun C1",
			ActorIndex: c.Index,
			AttackTag:  attacks.AttackTagNone,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       .5,
		}
		//3 blades
		for i := 0; i < 3; i++ {
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHit(
					c.Core.Combat.Player(),
					c.Core.Combat.PrimaryTarget(),
					nil,
					1,
				),
				attackHitmarks[c.NormalCounter]+i*5,
				attackHitmarks[c.NormalCounter]+i*5,
				c4CB,
			)
		}
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}
