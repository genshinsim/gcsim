package mika

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

const normalHitNum = 5

// based on raiden frames
// TODO: update frames, hitlags & hitboxes
var (
	attackFrames   [][]int
	attackHitmarks = [][]int{{14}, {9}, {14}, {14, 27}, {34}}
	// same between polearm and burst attacks so just use these arrays for both
	attackHitlagHaltFrame = [][]float64{{0.02}, {0.02}, {0.02}, {0, 0}, {0.04}}
	attackDefHalt         = [][]bool{{false}, {true}, {false}, {true, true}, {true}}
	attackHitboxes        = [][][]float64{{{1.8, 2.5}}, {{1.6}}, {{2.5, 2.5}}, {{4}, {4}}, {{1.5, 5}}}
	attackOffsets         = []float64{0, 0.2, 0.5, -0.3, 1.0}
	attackFanAngles       = []float64{360, 360, 360, 45, 360}
)

func init() {
	// NA cancels (polearm)
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 24)
	attackFrames[0][action.ActionAttack] = 18

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 26)
	attackFrames[1][action.ActionAttack] = 13

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 36)
	attackFrames[2][action.ActionAttack] = 26

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][1], 57)
	attackFrames[3][action.ActionAttack] = 41

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 50)
	attackFrames[4][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeSlash,
			Element:            attributes.Physical,
			Durability:         25,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			HitlagFactor:       0.01,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}
		ai.Mult = mult[c.TalentLvlAttack()]
		ap := combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][i][0],
			attackHitboxes[c.NormalCounter][i][1],
		)
		if c.NormalCounter == 1 {
			ap = combat.NewCircleHitOnTargetFanAngle(
				c.Core.Combat.Player(),
				geometry.Point{Y: attackOffsets[c.NormalCounter]},
				attackHitboxes[c.NormalCounter][i][0],
				attackFanAngles[c.NormalCounter],
			)
		} else if c.NormalCounter == 2 || c.NormalCounter == 3 {
			ai.StrikeType = attacks.StrikeTypeSpear
		}
		c.QueueCharTask(func() {
			c.Core.QueueAttack(ai, ap, 0, 0)
		}, attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:              frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength:     attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:       attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:               action.NormalAttackState,
		FramePausedOnHitlag: c.FramePausedOnHitlag,
	}
}
