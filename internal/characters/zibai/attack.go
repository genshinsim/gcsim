package zibai

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
	attackFrames               [][]int
	attackHitmarks             = [][]int{{11}, {12}, {17, 17 + 10}, {27}}
	attackHitlagHaltFrame      = [][]float64{{0.01}, {0.02}, {0.02, 0.01}, {0.03}}
	skillAttackHitlagHaltFrame = [][]float64{{0.01}, {0.02}, {0.02, 0.02}, {0.03}}
	attackDefHalt              = [][]bool{{true}, {true}, {true, false}, {true}}
	attackHitboxes             = [][]float64{{2}, {2, 3.3}, {3.6, 2.5}, {2.8}}
	attackOffsets              = [][]float64{{0.65}, {-0.1}, {0.2, -0.2}, {0.9}}
	attackPoiseDmg             = [][]float64{{52.9}, {48.7}, {32.3, 32.3}, {81.5}}
)

// 242
const normalHitNum = 4

func init() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 40)
	attackFrames[0][action.ActionAttack] = 18
	attackFrames[0][action.ActionCharge] = 21

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 36)
	attackFrames[1][action.ActionAttack] = 18
	attackFrames[1][action.ActionWalk] = 29

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 45)
	attackFrames[2][action.ActionAttack] = 40
	attackFrames[2][action.ActionWalk] = 40

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 67)
	attackFrames[3][action.ActionAttack] = 57
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillKey) {
		return c.skillAttack()
	}
	counter := c.NormalCounter
	for i, delay := range attackHitmarks[counter] {
		c.QueueCharTask(func() {
			ai := info.AttackInfo{
				ActorIndex:         c.Index(),
				Abil:               fmt.Sprintf("Normal %v", counter),
				AttackTag:          attacks.AttackTagNormal,
				ICDTag:             attacks.ICDTagNormalAttack,
				ICDGroup:           attacks.ICDGroupDefault,
				StrikeType:         attacks.StrikeTypeSlash,
				Element:            attributes.Physical,
				Durability:         25,
				Mult:               attack[counter][i][c.TalentLvlAttack()],
				HitlagFactor:       0.01,
				HitlagHaltFrames:   attackHitlagHaltFrame[counter][i] * 60,
				CanBeDefenseHalted: attackDefHalt[counter][i],
			}

			ap := combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				info.Point{Y: attackOffsets[counter][i]},
				attackHitboxes[counter][0],
			)

			switch counter {
			case 1, 2:
				ap = combat.NewBoxHitOnTarget(
					c.Core.Combat.Player(),
					info.Point{Y: attackOffsets[counter][i]},
					attackHitboxes[counter][0],
					attackHitboxes[counter][1],
				)
			}

			c.Core.QueueAttack(ai, ap, 0, 0)
		}, delay)
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[counter][action.InvalidAction],
		CanQueueAfter:   attackFrames[counter][action.ActionDash],
		State:           action.NormalAttackState,
	}, nil
}

func (c *char) skillAttack() (action.Info, error) {
	counter := c.NormalCounter
	for i, mult := range skillAttack[counter] {
		c.QueueCharTask(func() {
			ai := info.AttackInfo{
				ActorIndex:         c.Index(),
				Abil:               fmt.Sprintf("Normal %v", counter),
				AttackTag:          attacks.AttackTagNormal,
				ICDTag:             attacks.ICDTagNormalAttack,
				ICDGroup:           attacks.ICDGroupDefault,
				StrikeType:         attacks.StrikeTypeBlunt,
				PoiseDMG:           attackPoiseDmg[counter][i],
				Element:            attributes.Geo,
				IgnoreInfusion:     true,
				Durability:         25,
				Mult:               mult[c.TalentLvlAttack()],
				UseDef:             true,
				HitlagFactor:       0.01,
				HitlagHaltFrames:   skillAttackHitlagHaltFrame[counter][i] * 60,
				CanBeDefenseHalted: attackDefHalt[counter][i],
			}

			ap := combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				info.Point{Y: attackOffsets[counter][i]},
				attackHitboxes[counter][0],
			)

			switch counter {
			case 1, 2:
				ap = combat.NewBoxHitOnTarget(
					c.Core.Combat.Player(),
					info.Point{Y: attackOffsets[counter][i]},
					attackHitboxes[counter][0],
					attackHitboxes[counter][1],
				)
			}

			c.Core.QueueAttack(ai, ap, 0, 0)

			if counter == 3 && c.Core.Player.GetMoonsignLevel() >= 2 {
				c.skillLastAttack()
			}
		}, attackHitmarks[counter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[counter][action.InvalidAction],
		CanQueueAfter:   attackFrames[counter][action.ActionDash],
		State:           action.NormalAttackState,
	}, nil
}

func (c *char) skillLastAttack() {
	c.QueueCharTask(func() {
		ai := info.AttackInfo{
			ActorIndex:       c.Index(),
			Abil:             "Normal 4" + lunarCrystallizeAbil,
			AttackTag:        attacks.AttackTagDirectLunarCrystallize,
			ICDTag:           attacks.ICDTagNone,
			ICDGroup:         attacks.ICDGroupDefault,
			StrikeType:       attacks.StrikeTypeDefault,
			Element:          attributes.Geo,
			IgnoreInfusion:   true,
			Durability:       25,
			Mult:             skillLastAttackBonus[c.TalentLvlAttack()] * c.c4N4Bonus(),
			UseDef:           true,
			IgnoreDefPercent: 1,
		}

		ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: 2.5}, 4)
		c.Core.QueueAttack(ai, ap, 0, 0, c.particleCB, c.radianceCB)
	}, 20)
}
