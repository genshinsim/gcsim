package varka

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
	attackFrames          [][]int
	attackHitmarks        = [][]int{{21}, {14, 26}, {28, 33}, {28, 33}, {30, 30}}
	attackPoiseDMG        = []float64{107.0, 48.8, 44.43, 49.9, 155.4}
	attackHitlagHaltFrame = []float64{0.05, 0.05, 0.05, 0.05, 0.1}
	attackHitboxes        = [][]float64{{2.2}, {3.3, 4.3}, {2.8, 5.0}, {2.8, 5.0}, {3.2}}
	attackOffsets         = []float64{0.5, -1.3, 0.5, 0.5, -0.8}

	skillAttackFrames          [][]int
	skillAttackHitmarks        = [][]int{{21}, {14, 26}, {28, 33}, {28, 33}, {30, 30}}
	skillAttackPoiseDMG        = []float64{107.0, 48.8, 44.43, 49.9, 155.4}
	skillAttackHitlagHaltFrame = []float64{0.05, 0.05, 0.05, 0.05, 0.1}
	skillAttackHitboxes        = [][]float64{{2.2}, {3.3, 4.3}, {2.8, 5.0}, {2.8, 5.0}, {3.2}}
	skillAttackOffsets         = []float64{0.5, -1.3, 0.5, 0.5, -0.8}
)

const (
	normalHitNum = 5
)

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 40) // N1 -> Walk
	attackFrames[0][action.ActionAttack] = 31                                // N1 -> N2
	attackFrames[0][action.ActionCharge] = 31                                // N1 -> CA

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][1], 50) // N2 -> Walk
	attackFrames[1][action.ActionAttack] = 42                                // N2 -> N3
	attackFrames[1][action.ActionCharge] = 42                                // N2 -> CA

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 59) // N3 -> Walk
	attackFrames[2][action.ActionAttack] = 46                                // N3 -> N4
	attackFrames[2][action.ActionCharge] = 47                                // N3 -> CA

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][1], 59) // N3 -> Walk
	attackFrames[3][action.ActionAttack] = 46                                // N3 -> N4
	attackFrames[3][action.ActionCharge] = 47                                // N3 -> CA

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][1], 61) // N4 -> Walk
	attackFrames[4][action.ActionAttack] = 60                                // N4 -> N1
	attackFrames[4][action.ActionCharge] = 61                                // N4 -> CA

	skillAttackFrames = make([][]int, normalHitNum)
	skillAttackFrames[0] = frames.InitNormalCancelSlice(skillAttackHitmarks[0][0], 40) // N1 -> Walk
	skillAttackFrames[0][action.ActionAttack] = 31                                     // N1 -> N2
	skillAttackFrames[0][action.ActionCharge] = 31                                     // N1 -> CA

	skillAttackFrames[1] = frames.InitNormalCancelSlice(skillAttackHitmarks[1][1], 50) // N2 -> Walk
	skillAttackFrames[1][action.ActionAttack] = 42                                     // N2 -> N3
	skillAttackFrames[1][action.ActionCharge] = 42                                     // N2 -> CA

	skillAttackFrames[2] = frames.InitNormalCancelSlice(skillAttackHitmarks[2][1], 59) // N3 -> Walk
	skillAttackFrames[2][action.ActionAttack] = 46                                     // N3 -> N4
	skillAttackFrames[2][action.ActionCharge] = 47                                     // N3 -> CA

	skillAttackFrames[3] = frames.InitNormalCancelSlice(skillAttackHitmarks[3][1], 59) // N3 -> Walk
	skillAttackFrames[3][action.ActionAttack] = 46                                     // N3 -> N4
	skillAttackFrames[3][action.ActionCharge] = 47                                     // N3 -> CA

	skillAttackFrames[4] = frames.InitNormalCancelSlice(skillAttackHitmarks[4][1], 61) // N4 -> Walk
	skillAttackFrames[4][action.ActionAttack] = 60                                     // N4 -> N1
	skillAttackFrames[4][action.ActionCharge] = 61                                     // N4 -> CA
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillKey) {
		return c.skillAttack()
	}
	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		info.Point{Y: attackOffsets[c.NormalCounter]},
		attackHitboxes[c.NormalCounter][0],
	)

	for i, delay := range attackHitmarks[c.NormalCounter] {
		ai := info.AttackInfo{
			ActorIndex:         c.Index(),
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeBlunt,
			PoiseDMG:           attackPoiseDMG[c.NormalCounter],
			Element:            attributes.Physical,
			Durability:         25,
			Mult:               attack[c.NormalCounter][i][c.TalentLvlAttack()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter] * 60,
			CanBeDefenseHalted: true,
		}
		c.Core.QueueAttack(ai, ap, delay, delay)
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}, nil
}

func (c *char) skillAttack() (action.Info, error) {
	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		info.Point{Y: skillAttackOffsets[c.NormalCounter]},
		skillAttackHitboxes[c.NormalCounter][0],
	)

	ele := []attributes.Element{c.conversionElem, attributes.Anemo}
	offset := 0
	switch c.NormalCounter {
	case 1, 2:
		offset = 1
	}

	for i, delay := range skillAttackHitmarks[c.NormalCounter] {
		ai := info.AttackInfo{
			ActorIndex:         c.Index(),
			Abil:               fmt.Sprintf("Sturm und Drang Normal %v", c.NormalCounter),
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeBlunt,
			PoiseDMG:           skillAttackPoiseDMG[c.NormalCounter],
			Element:            ele[(i+offset)%2],
			Durability:         25,
			Mult:               skillAttack[c.NormalCounter][i][c.TalentLvlSkill()] * c.a1SkillMulti(),
			HitlagFactor:       0.01,
			HitlagHaltFrames:   skillAttackHitlagHaltFrame[c.NormalCounter] * 60,
			CanBeDefenseHalted: true,
		}
		c.Core.QueueAttack(ai, ap, delay, delay, c.fourWindsCDRedCB)
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}, nil
}
