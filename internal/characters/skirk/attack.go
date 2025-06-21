package skirk

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
	attackHitmarks        = [][]int{{12}, {10}, {7, 12}, {10}, {25}}
	attackHitlagHaltFrame = [][]float64{{0.03}, {0.03}, {0.03, 0.03}, {0.03}, {0.12}}
	attackHitlagFactor    = [][]float64{{0.01}, {0.01}, {0.05, 0.05}, {0.05}, {0.01}}
	attackHitboxes        = [][]float64{{1.2}, {1.4, 2.2}, {1.6}, {1.6}, {2.2}}
	attackOffsets         = [][]float64{{0.8}, {0}, {1, 0.6}, {0.6, 0.6}, {1}}
	attackFanAngles       = [][]float64{{360}, {360}, {30, 360}, {360}, {360}}

	attackSkillFrames          [][]int
	attackSkillHitmarks        = [][]int{{15}, {8}, {11, 22}, {12, 29}, {24}}
	attackSkillHitlagHaltFrame = [][]float64{{0.01}, {0.01}, {0.01, 0.00}, {0.00, 0.0}, {0.05}}
	attackSkillHitlagFactor    = [][]float64{{0.05}, {0.05}, {0.05, 0.01}, {0.05, 0.05}, {0.05}}
	attackSkillHitboxes        = [][]float64{{1.2}, {1.4, 2.2}, {1.6}, {1.6}, {2.2}}
	attackSkillOffsets         = [][]float64{{0.8}, {0}, {1, 0.6}, {0.6, 0.6}, {1}}
	attackSkillFanAngles       = [][]float64{{360}, {360}, {30, 360}, {360, 360}, {360}}
)

const normalHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 46) // N1 -> W
	attackFrames[0][action.ActionCharge] = 23                                // N1 -> CA
	attackFrames[0][action.ActionAttack] = 14                                // N1 -> N2

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 40) // N2 -> W
	attackFrames[1][action.ActionAttack] = 21                                // N2 -> N3
	attackFrames[1][action.ActionCharge] = 20                                // N2 -> CA

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 74) // N3 -> W
	attackFrames[2][action.ActionAttack] = 37                                // N3 -> N4
	attackFrames[2][action.ActionCharge] = 37                                // N3 -> CA

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 40) // N4 -> W
	attackFrames[3][action.ActionAttack] = 20                                // N4 -> N5
	attackFrames[3][action.ActionCharge] = 26                                // N4 -> CA

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 90) // N5 -> W
	attackFrames[4][action.ActionAttack] = 45                                // N4 -> N5
	attackFrames[4][action.ActionCharge] = 45                                // N4 -> CA

	attackSkillFrames = make([][]int, normalHitNum)

	attackSkillFrames[0] = frames.InitNormalCancelSlice(attackSkillHitmarks[0][0], 240) // N1 -> W
	attackSkillFrames[0][action.ActionCharge] = 18                                      // N1 -> CA
	attackSkillFrames[0][action.ActionAttack] = 18                                      // N1 -> N2

	attackSkillFrames[1] = frames.InitNormalCancelSlice(attackSkillHitmarks[1][0], 240) // N2 -> W
	attackSkillFrames[1][action.ActionAttack] = 19                                      // N2 -> N3
	attackSkillFrames[1][action.ActionCharge] = 19                                      // N2 -> CA

	attackSkillFrames[2] = frames.InitNormalCancelSlice(attackSkillHitmarks[2][1], 240) // N3 -> W
	attackSkillFrames[2][action.ActionAttack] = 34                                      // N3 -> N4
	attackSkillFrames[2][action.ActionCharge] = 34                                      // N3 -> CA

	attackSkillFrames[3] = frames.InitNormalCancelSlice(attackSkillHitmarks[3][1], 240) // N4 -> W
	attackSkillFrames[3][action.ActionAttack] = 33                                      // N4 -> N5
	attackSkillFrames[3][action.ActionCharge] = 33                                      // N4 -> CA

	attackSkillFrames[4] = frames.InitNormalCancelSlice(attackSkillHitmarks[4][0], 240) // N5 -> W
	attackSkillFrames[4][action.ActionAttack] = 54                                      // N5 -> N1
	attackSkillFrames[4][action.ActionCharge] = 54                                      // N5 -> CA
}

// Standard attack - nothing special
func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillKey) {
		return c.AttackSkill(p)
	}

	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			Mult:               mult[c.TalentLvlAttack()],
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeSlash,
			Element:            attributes.Physical,
			Durability:         25,
			HitlagFactor:       attackHitlagFactor[c.NormalCounter][i],
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: true,
		}
		ap := combat.NewCircleHitOnTargetFanAngle(
			c.Core.Combat.Player(),
			geometry.Point{Y: attackOffsets[c.NormalCounter][i]},
			attackHitboxes[c.NormalCounter][0],
			attackFanAngles[c.NormalCounter][i],
		)
		if c.NormalCounter == 1 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: attackOffsets[c.NormalCounter][i]},
				attackHitboxes[c.NormalCounter][0],
				attackHitboxes[c.NormalCounter][1],
			)
		}
		c.QueueCharTask(func() {
			c.Core.QueueAttack(ai, ap, 0, 0)
		}, attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}, nil
}

// Standard attack - nothing special
func (c *char) AttackSkill(p map[string]int) (action.Info, error) {
	for i, mult := range skillAttack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:       c.Index,
			Abil:             fmt.Sprintf("Normal (Skill) %v", c.NormalCounter),
			Mult:             mult[c.TalentLvlSkill()] * c.a4MultAttack(),
			AttackTag:        attacks.AttackTagNormal,
			ICDTag:           attacks.ICDTagNormalAttack,
			ICDGroup:         attacks.ICDGroupDefault,
			StrikeType:       attacks.StrikeTypeSlash,
			Element:          attributes.Cryo,
			Durability:       25,
			HitlagFactor:     attackSkillHitlagFactor[c.NormalCounter][i],
			HitlagHaltFrames: attackSkillHitlagHaltFrame[c.NormalCounter][i] * 60,
			IgnoreInfusion:   true,
		}
		ap := combat.NewCircleHitOnTargetFanAngle(
			c.Core.Combat.Player(),
			geometry.Point{Y: attackSkillOffsets[c.NormalCounter][i]},
			attackSkillHitboxes[c.NormalCounter][0],
			attackSkillFanAngles[c.NormalCounter][i],
		)
		if c.NormalCounter == 1 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: attackSkillOffsets[c.NormalCounter][i]},
				attackSkillHitboxes[c.NormalCounter][0],
				attackSkillHitboxes[c.NormalCounter][1],
			)
		}
		c6cb := c.c6OnAttackCB()
		c.QueueCharTask(func() {
			c.Core.QueueAttack(ai, ap, 0, 0, c6cb)
		}, attackSkillHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackSkillFrames),
		AnimationLength: attackSkillFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackSkillHitmarks[c.NormalCounter][len(attackSkillHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}, nil
}
