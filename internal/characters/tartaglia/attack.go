package tartaglia

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

const normalHitNum = 6

var (
	attackFrames   [][]int
	attackHitmarks = []int{17, 8, 15, 19, 11, 14}
)

func init() {
	// attack (ranged) -> x
	attackFrames = make([][]int, normalHitNum)

	// N1 -> x
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 26)

	// N2 -> x
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 27)

	// N3 -> x
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 33)

	// N4 -> x
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 32)

	// N5 -> x
	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4], 33)

	// N6 -> x
	attackFrames[5] = frames.InitNormalCancelSlice(attackHitmarks[5], 66)
}

// Normal attack
// Perform up to 6 consecutive shots with a bow.
func (c *char) Attack(p map[string]int) action.ActionInfo {
	if c.StatusIsActive(MeleeKey) {
		return c.meleeAttack(p)
	}

	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  combat.AttackTagNormal,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypePierce,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewBoxHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			combat.Point{Y: -0.5},
			0.1,
			1,
		),
		attackHitmarks[c.NormalCounter],
		attackHitmarks[c.NormalCounter]+travel,
	)

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}

var (
	meleeFrames           [][]int
	meleeHitmarks         = [][]int{{8}, {6}, {16}, {7}, {7}, {4, 20}}
	meleeHitlagHaltFrames = [][]float64{{0.03}, {0.03}, {0.06}, {0.06}, {0.06}, {0.03, 0.12}}
	meleeHitboxes         = [][][]float64{{{1.8}}, {{1.8}}, {{2}}, {{2}}, {{2.2}}, {{2, 3.2}, {2.2}}}
	meleeOffsets          = [][]float64{{0.8}, {0.8}, {0.6}, {0.9}, {0.6}, {0.3, 1.5}}
	meleeFanAngles        = []float64{300, 270, 300, 300, 360, 360}
)

func init() {
	// attack (melee) -> x
	meleeFrames = make([][]int, normalHitNum)

	// N1 -> x
	meleeFrames[0] = frames.InitNormalCancelSlice(meleeHitmarks[0][0], 23)
	meleeFrames[0][action.ActionAttack] = 10
	meleeFrames[0][action.ActionCharge] = 23

	// N2 -> x
	meleeFrames[1] = frames.InitNormalCancelSlice(meleeHitmarks[1][0], 23)
	meleeFrames[1][action.ActionAttack] = 11
	meleeFrames[1][action.ActionCharge] = 23

	// N3 -> x
	meleeFrames[2] = frames.InitNormalCancelSlice(meleeHitmarks[2][0], 37)
	meleeFrames[2][action.ActionAttack] = 32
	meleeFrames[2][action.ActionCharge] = 37

	// N4 -> x
	meleeFrames[3] = frames.InitNormalCancelSlice(meleeHitmarks[3][0], 37)
	meleeFrames[3][action.ActionAttack] = 33
	meleeFrames[3][action.ActionCharge] = 37

	// N5 -> x
	meleeFrames[4] = frames.InitNormalCancelSlice(meleeHitmarks[4][0], 23)
	meleeFrames[4][action.ActionAttack] = 22
	meleeFrames[4][action.ActionCharge] = 23

	// N6 -> x
	meleeFrames[5] = frames.InitNormalCancelSlice(meleeHitmarks[5][1], 65)
	meleeFrames[5][action.ActionAttack] = 65
	meleeFrames[5][action.ActionCharge] = 500 // illegal action
}

// Melee stance attack.
// Perform up to 6 consecutive Hydro strikes.
func (c *char) meleeAttack(p map[string]int) action.ActionInfo {
	for i, mult := range eAttack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:          combat.AttackTagNormal,
			ICDTag:             combat.ICDTagNormalAttack,
			ICDGroup:           combat.ICDGroupDefault,
			StrikeType:         combat.StrikeTypeSlash,
			Element:            attributes.Hydro,
			Durability:         25,
			HitlagFactor:       0.01,
			CanBeDefenseHalted: true,
			Mult:               mult[c.TalentLvlSkill()],
			HitlagHaltFrames:   meleeHitlagHaltFrames[c.NormalCounter][i] * 60,
		}
		ap := combat.NewCircleHitOnTargetFanAngle(
			c.Core.Combat.Player(),
			combat.Point{Y: meleeOffsets[c.NormalCounter][i]},
			meleeHitboxes[c.NormalCounter][i][0],
			meleeFanAngles[c.NormalCounter],
		)
		if c.NormalCounter == 5 && i == 0 {
			ai.StrikeType = combat.StrikeTypeSpear
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				combat.Point{Y: meleeOffsets[c.NormalCounter][i]},
				meleeHitboxes[c.NormalCounter][i][0],
				meleeHitboxes[c.NormalCounter][i][1],
			)
		}
		c.QueueCharTask(func() {
			c.Core.QueueAttack(
				ai,
				ap,
				0,
				0,
				c.meleeApplyRiptide, // riptide can trigger on the same hit that applies
				c.rtSlashCallback,
			)
		}, meleeHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, meleeFrames),
		AnimationLength: meleeFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   meleeHitmarks[c.NormalCounter][len(meleeHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}
