package skirk

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

const windup = 2

var (
	attackFrames [][]int

	// The n1 hitmark includes 2f of windup
	attackHitmarks        = [][]int{{11 + windup}, {7}, {8, 8 + 14}, {11}, {35}}
	attackHitlagHaltFrame = [][]float64{{0.02}, {0.03}, {0.03, 0.00}, {0.05}, {0.06}}
	attackHitlagFactor    = [][]float64{{0.01}, {0.01}, {0.01, 0.00}, {0.01}, {0.01}}
	attackHitboxes        = [][][]float64{{{2}}, {{2.2}}, {{2.5, 3.2}, {2.5, 3.2}}, {{2.5, 2.4}}, {{3.2}}}
	attackOffsets         = [][]float64{{-0.1}, {-0.1}, {1.5, 1.5}, {1.1}, {-0.1}}

	attackSkillFrames          [][]int
	attackSkillHitmarks        = [][]int{{10 + windup}, {11}, {11, 11 + 12}, {11, 11 + 16}, {25}}
	attackSkillHitlagHaltFrame = [][]float64{{0.02}, {0.03}, {0.03, 0.00}, {0.03, 0.00}, {0.06}}
	attackSkillHitlagFactor    = [][]float64{{0.01}, {0.01}, {0.01, 0.00}, {0.01, 0.00}, {0.01}}
	attackSkillHitboxes        = [][][]float64{{{7, 2.4}}, {{7, 3}}, {{6, 3.6}, {6, 3.6}}, {{5, 3.6}, {5, 3.6}}, {{11, 3.6}}}
	attackSkillOffsets         = [][]float64{{0.9}, {1.4}, {1.7, 1.7}, {1.7, 1.7}, {1.7}}
)

const normalHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 25+windup) // N1 -> W
	attackFrames[0][action.ActionCharge] = 18 + windup                              // N1 -> CA
	attackFrames[0][action.ActionAttack] = 17 + windup                              // N1 -> N2

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 25) // N2 -> W
	attackFrames[1][action.ActionAttack] = 18                                // N2 -> N3
	attackFrames[1][action.ActionCharge] = 21                                // N2 -> CA

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 43) // N3 -> W
	attackFrames[2][action.ActionAttack] = 37                                // N3 -> N4
	attackFrames[2][action.ActionCharge] = 30                                // N3 -> CA

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 23) // N4 -> W
	attackFrames[3][action.ActionAttack] = 19                                // N4 -> N5
	attackFrames[3][action.ActionCharge] = 20                                // N4 -> CA

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 72) // N5 -> W
	attackFrames[4][action.ActionAttack] = 67                                // N4 -> N5
	attackFrames[4][action.ActionCharge] = 48                                // N4 -> CA

	attackSkillFrames = make([][]int, normalHitNum)

	attackSkillFrames[0] = frames.InitNormalCancelSlice(attackSkillHitmarks[0][0], 30) // N1 -> W
	attackSkillFrames[0][action.ActionAttack] = 10 + windup                            // N1 -> N2
	attackSkillFrames[0][action.ActionCharge] = 11 + windup                            // N1 -> CA

	attackSkillFrames[1] = frames.InitNormalCancelSlice(attackSkillHitmarks[1][0], 43) // N2 -> W
	attackSkillFrames[1][action.ActionAttack] = 21                                     // N2 -> N3
	attackSkillFrames[1][action.ActionCharge] = 12                                     // N2 -> CA

	attackSkillFrames[2] = frames.InitNormalCancelSlice(attackSkillHitmarks[2][1], 42) // N3 -> W
	attackSkillFrames[2][action.ActionAttack] = 31                                     // N3 -> N4
	attackSkillFrames[2][action.ActionCharge] = 34                                     // N3 -> CA

	attackSkillFrames[3] = frames.InitNormalCancelSlice(attackSkillHitmarks[3][1], 60) // N4 -> W
	attackSkillFrames[3][action.ActionAttack] = 33                                     // N4 -> N5
	attackSkillFrames[3][action.ActionCharge] = 40                                     // N4 -> CA

	attackSkillFrames[4] = frames.InitNormalCancelSlice(attackSkillHitmarks[4][0], 72) // N5 -> W
	attackSkillFrames[4][action.ActionAttack] = 51                                     // N5 -> N1
	attackSkillFrames[4][action.ActionCharge] = 42                                     // N5 -> CA
}

// Standard attack - nothing special
func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillKey) {
		return c.AttackSkill(p)
	}

	windupRemove := 0

	// skip N1 windup out of NA and Q
	if c.NormalCounter == 0 {
		switch c.Core.Player.CurrentState() {
		case action.BurstState:
			windupRemove = windup
		case action.NormalAttackState:
			windupRemove = windup
		}
	}

	for i, mult := range attack[c.NormalCounter] {
		ai := info.AttackInfo{
			ActorIndex:       c.Index(),
			Abil:             fmt.Sprintf("Normal %v", c.NormalCounter),
			Mult:             mult[c.TalentLvlAttack()],
			AttackTag:        attacks.AttackTagNormal,
			ICDTag:           attacks.ICDTagNormalAttack,
			ICDGroup:         attacks.ICDGroupDefault,
			StrikeType:       attacks.StrikeTypeSlash,
			Element:          attributes.Physical,
			Durability:       25,
			HitlagFactor:     attackHitlagFactor[c.NormalCounter][i],
			HitlagHaltFrames: attackHitlagHaltFrame[c.NormalCounter][i] * 60,
		}
		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			info.Point{Y: attackOffsets[c.NormalCounter][i]},
			attackHitboxes[c.NormalCounter][i][0],
		)

		if i == 3 {
			ai.StrikeType = attacks.StrikeTypeSpear
		}

		if c.NormalCounter == 2 || c.NormalCounter == 3 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				info.Point{Y: attackOffsets[c.NormalCounter][i]},
				attackHitboxes[c.NormalCounter][i][0],
				attackHitboxes[c.NormalCounter][i][1],
			)
		}
		c.QueueCharTask(func() {
			c.Core.QueueAttack(ai, ap, 0, 0)
		}, attackHitmarks[c.NormalCounter][i]-windup)
	}

	c.prevNASkillState = false
	normalCounter := c.NormalCounter

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames: func(next action.Action) int {
			return frames.AtkSpdAdjust(attackFrames[normalCounter][next], c.Stat(attributes.AtkSpd)) - windupRemove
		},
		AnimationLength: attackFrames[normalCounter][action.InvalidAction] - windupRemove,
		CanQueueAfter:   attackHitmarks[normalCounter][len(attackHitmarks[normalCounter])-1] - windupRemove,
		State:           action.NormalAttackState,
	}, nil
}

// Standard attack - nothing special
func (c *char) AttackSkill(p map[string]int) (action.Info, error) {
	windupRemove := 0

	// skip N1 windup out of NA
	if c.NormalCounter == 0 && c.Core.Player.CurrentState() == action.NormalAttackState {
		windupRemove = windup
	}

	for i, mult := range skillAttack[c.NormalCounter] {
		ai := info.AttackInfo{
			ActorIndex:       c.Index(),
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

		if i == 3 {
			ai.StrikeType = attacks.StrikeTypeSpear
		}

		ap := combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			info.Point{Y: attackSkillOffsets[c.NormalCounter][i]},
			attackSkillHitboxes[c.NormalCounter][i][0],
			attackSkillHitboxes[c.NormalCounter][i][1],
		)

		c6cb := c.c6OnAttackCB()
		c.QueueCharTask(func() {
			c.Core.QueueAttack(ai, ap, 0, 0, c6cb)
		}, attackSkillHitmarks[c.NormalCounter][i]-windupRemove)
	}

	c.prevNASkillState = true
	normalCounter := c.NormalCounter
	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames: func(next action.Action) int {
			return frames.AtkSpdAdjust(attackSkillFrames[normalCounter][next], c.Stat(attributes.AtkSpd)) - windupRemove
		},
		AnimationLength: attackSkillFrames[normalCounter][action.InvalidAction] - windupRemove,
		CanQueueAfter:   attackSkillHitmarks[normalCounter][len(attackSkillHitmarks[normalCounter])-1] - windupRemove,
		State:           action.NormalAttackState,
	}, nil
}
