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
	attackFrames             [][]int
	attackHitmarks           = [][]int{{19}, {18, 18 + 10}, {27, 27 + 16}, {19, 19 + 5}, {44, 44 + 1}}
	attackPoiseDMG           = [][]float64{{87.5334}, {32.0776, 59.5728}, {43.3742, 80.5521}, {74.1236, 39.9127}, {93.2701, 50.2223}}
	attackHitlagHaltFrame    = [][]float64{{0.03}, {0, 0.06}, {0, 0.06}, {0, 0.09}, {0, 0.1}}
	attackCanBeDefenseHalted = [][]bool{{true}, {false, true}, {false, true}, {false, true}, {false, true}}
	attackHitboxes           = [][]float64{{2, 3.2}, {2.5}, {3}, {2.8}, {2.8}}
	attackOffsets            = [][]float64{{-0.5}, {0.5, 0.5}, {1, 1.5}, {1.2, 1.2}, {0.6, 0.6}}
)

const (
	normalHitNum = 5
)

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 46) // N1 -> Walk
	attackFrames[0][action.ActionAttack] = 23                                // N1 -> N2
	attackFrames[0][action.ActionCharge] = 22                                // N1 -> CA

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][1], 46) // N2 -> Walk
	attackFrames[1][action.ActionAttack] = 29                                // N2 -> N3
	attackFrames[1][action.ActionCharge] = 30                                // N2 -> CA

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 60) // N3 -> Walk
	attackFrames[2][action.ActionAttack] = 55                                // N3 -> N4
	attackFrames[2][action.ActionCharge] = 48                                // N3 -> CA

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][1], 47) // N3 -> Walk
	attackFrames[3][action.ActionAttack] = 40                                // N3 -> N4
	attackFrames[3][action.ActionCharge] = 28                                // N3 -> CA

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][1], 82) // N4 -> Walk
	attackFrames[4][action.ActionAttack] = 73                                // N4 -> N1
	attackFrames[4][action.ActionCharge] = 48                                // N4 -> CA
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillKey) {
		return c.skillAttack()
	}

	windup := 0
	switch c.Core.Player.CurrentState() {
	case action.Idle, action.SwapState, action.JumpState, action.DashState:
		windup = 5
	}

	for i, delay := range attackHitmarks[c.NormalCounter] {
		ai := info.AttackInfo{
			ActorIndex:         c.Index(),
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeBlunt,
			PoiseDMG:           attackPoiseDMG[c.NormalCounter][i],
			Element:            attributes.Physical,
			Durability:         25,
			Mult:               attack[c.NormalCounter][i][c.TalentLvlAttack()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackCanBeDefenseHalted[c.NormalCounter][i],
		}
		var ap info.AttackPattern
		hitbox := attackHitboxes[c.NormalCounter]
		switch len(hitbox) {
		case 1:
			ap = combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				info.Point{Y: attackOffsets[c.NormalCounter][i]},
				hitbox[0],
			)
		case 2:
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				info.Point{Y: attackOffsets[c.NormalCounter][i]},
				hitbox[0],
				hitbox[1],
			)
		default:
			panic("varka NA hitbox array incorrect size")
		}
		c.Core.QueueAttack(ai, ap, delay+windup, delay+windup)
	}

	normalCounter := c.NormalCounter

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFuncWithOffset(c.Character, attackFrames, windup),
		AnimationLength: attackFrames[normalCounter][action.InvalidAction] + windup,
		CanQueueAfter:   attackHitmarks[normalCounter][len(attackHitmarks[c.NormalCounter])-1] + windup,
		State:           action.NormalAttackState,
	}, nil
}

func (c *char) skillAttack() (action.Info, error) {
	ele := []attributes.Element{c.conversionElem, attributes.Anemo}
	offset := 0
	switch c.NormalCounter {
	case 1, 2:
		offset = 1
	}

	windup := 0
	switch c.Core.Player.CurrentState() {
	case action.Idle, action.SwapState, action.JumpState, action.DashState, action.WalkState:
		windup = 5
	}

	for i, delay := range attackHitmarks[c.NormalCounter] {
		ai := info.AttackInfo{
			ActorIndex:         c.Index(),
			Abil:               fmt.Sprintf("Sturm und Drang Normal %v", c.NormalCounter),
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeBlunt,
			PoiseDMG:           attackPoiseDMG[c.NormalCounter][i],
			Element:            ele[(i+offset)%2],
			Durability:         25,
			Mult:               skillAttack[c.NormalCounter][i][c.TalentLvlSkill()] * c.a1SkillMulti(),
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackCanBeDefenseHalted[c.NormalCounter][i],
		}
		var ap info.AttackPattern
		hitbox := attackHitboxes[c.NormalCounter]
		switch len(hitbox) {
		case 1:
			ap = combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				info.Point{Y: attackOffsets[c.NormalCounter][i]},
				hitbox[0],
			)
		case 2:
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				info.Point{Y: attackOffsets[c.NormalCounter][i]},
				hitbox[0],
				hitbox[1],
			)
		default:
			panic("varka NA hitbox array incorrect size")
		}

		c.Core.QueueAttack(ai, ap, delay+windup, delay+windup, c.fourWindsCDRedCB())
	}

	normalCounter := c.NormalCounter

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFuncWithOffset(c.Character, attackFrames, windup),
		AnimationLength: attackFrames[normalCounter][action.InvalidAction] + windup,
		CanQueueAfter:   attackHitmarks[normalCounter][len(attackHitmarks[c.NormalCounter])-1] + windup,
		State:           action.NormalAttackState,
	}, nil
}
