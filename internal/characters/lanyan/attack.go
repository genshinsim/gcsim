package lanyan

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
	attackFrames [][]int

	attackHitmarks        = [][]int{{11}, {17, 37}, {15, 21}, {40}}
	attackHitlagHaltFrame = [][]float64{{0.03}, {0.03, 0}, {0.03, 0}, {0.06}}
	attackDefHalt         = [][]bool{{true}, {true, true}, {false, false}, {true}}
	attackOffsets         = []float64{-0.2, -0.2, 0.3, 0.5}
	attackHitboxes        = [][]float64{{2.2, 3.0}, {2.3, 3.0}, {2.2}, {2.4}}
	attackFanAngles       = []float64{0, 0, 0, 240}
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 30) // N1 -> Walk
	attackFrames[0][action.ActionAttack] = 20
	attackFrames[0][action.ActionCharge] = 21

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][1], 46) // N2 -> N3
	attackFrames[1][action.ActionCharge] = 30
	attackFrames[1][action.ActionWalk] = 45

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 53) // N3 -> N4
	attackFrames[2][action.ActionCharge] = 37
	attackFrames[2][action.ActionWalk] = 47

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 63) // N4 -> Walk/N1
	attackFrames[3][action.ActionCharge] = 50
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	for i := 0; i < len(attack[c.NormalCounter]); i++ {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeDefault,
			Element:            attributes.Anemo,
			Durability:         25,
			Mult:               attack[c.NormalCounter][i][c.TalentLvlAttack()],
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}
		if ai.HitlagHaltFrames > 0 {
			ai.HitlagFactor = 0.01
		}

		var ap combat.AttackPattern
		switch {
		case len(attackHitboxes[c.NormalCounter]) == 2: // box
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: attackOffsets[c.NormalCounter]},
				attackHitboxes[c.NormalCounter][0],
				attackHitboxes[c.NormalCounter][1],
			)
		case attackFanAngles[c.NormalCounter] > 0: // circle with fan angle
			ap = combat.NewCircleHitOnTargetFanAngle(
				c.Core.Combat.Player(),
				geometry.Point{Y: attackOffsets[c.NormalCounter]},
				attackHitboxes[c.NormalCounter][0],
				attackFanAngles[c.NormalCounter],
			)
		default: // circle
			ap = combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: attackOffsets[c.NormalCounter]},
				attackHitboxes[c.NormalCounter][0],
			)
		}

		c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter][i], attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}, nil
}
