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
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], attackHitmarks[0])
	attackFrames[0][action.ActionAttack] = 26
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], attackHitmarks[1])
	attackFrames[1][action.ActionAttack] = 27
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], attackHitmarks[2])
	attackFrames[2][action.ActionAttack] = 33
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], attackHitmarks[3])
	attackFrames[3][action.ActionAttack] = 32
	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4], attackHitmarks[4])
	attackFrames[4][action.ActionAttack] = 33
	attackFrames[5] = frames.InitNormalCancelSlice(attackHitmarks[5], attackHitmarks[5])
	attackFrames[5][action.ActionAttack] = 66
}

// Normal attack
// Perform up to 6 consecutive shots with a bow.
func (c *char) Attack(p map[string]int) action.ActionInfo {
	if c.Core.Status.Duration("tartagliamelee") > 0 {
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
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypePierce,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewDefSingleTarget(c.Core.Combat.DefaultTarget, combat.TargettableEnemy),
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
	meleeHitmarks         = [][]int{{8}, {6}, {16}, {7}, {7}, {4, 16}}
	meleeHitlagHaltFrames = [][]float64{{0.03}, {0.03}, {0.06}, {0.06}, {0.06}, {0.03, 0.12}}
)

func init() {
	// attack (melee) -> x
	meleeFrames = make([][]int, normalHitNum)

	meleeFrames[0] = frames.InitNormalCancelSlice(meleeHitmarks[0][0], 23)
	meleeFrames[0][action.ActionAttack] = 10
	meleeFrames[0][action.ActionCharge] = 23
	meleeFrames[1] = frames.InitNormalCancelSlice(meleeHitmarks[1][0], 23)
	meleeFrames[1][action.ActionAttack] = 11
	meleeFrames[1][action.ActionCharge] = 23
	meleeFrames[2] = frames.InitNormalCancelSlice(meleeHitmarks[2][0], 37)
	meleeFrames[2][action.ActionAttack] = 32
	meleeFrames[2][action.ActionCharge] = 37
	meleeFrames[3] = frames.InitNormalCancelSlice(meleeHitmarks[3][0], 37)
	meleeFrames[3][action.ActionAttack] = 33
	meleeFrames[3][action.ActionCharge] = 37
	meleeFrames[4] = frames.InitNormalCancelSlice(meleeHitmarks[4][0], 23)
	meleeFrames[4][action.ActionAttack] = 22
	meleeFrames[4][action.ActionCharge] = 23
	meleeFrames[5] = frames.InitNormalCancelSlice(meleeHitmarks[5][0]+meleeHitmarks[5][1], 65)
	meleeFrames[5][action.ActionAttack] = 65
}

// Melee stance attack.
// Perform up to 6 consecutive Hydro strikes.
func (c *char) meleeAttack(p map[string]int) action.ActionInfo {
	lastMultiHit := 0
	for i, mult := range eAttack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:       c.Index,
			Abil:             fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:        combat.AttackTagNormal,
			ICDTag:           combat.ICDTagNormalAttack,
			ICDGroup:         combat.ICDGroupDefault,
			StrikeType:       combat.StrikeTypeSlash,
			Element:          attributes.Hydro,
			Durability:       25,
			HitlagFactor:     0.01,
			HitlagHaltFrames: meleeHitlagHaltFrames[c.NormalCounter][i] * 60,
			Mult:             mult[c.TalentLvlSkill()],
		}
		hitmark := lastMultiHit + meleeHitmarks[c.NormalCounter][i]
		c.QueueCharTask(func() {
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHit(c.Core.Combat.Player(), .5, false, combat.TargettableEnemy),
				0,
				0,
				c.meleeApplyRiptide, // riptide can trigger on the same hit that applies
				c.rtSlashCallback,
			)
		}, hitmark)
		lastMultiHit = hitmark
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, meleeFrames),
		AnimationLength: meleeFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   meleeHitmarks[c.NormalCounter][len(meleeHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}
