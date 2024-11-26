package mavuika

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
	// imaginary numbers
	attackFrames          [][]int
	attackHitmarks        = [][]int{{21}, {22, 22 + 9}, {34, 34 + 9, 34 + 9 + 9}, {42}}
	attackHitlagHaltFrame = [][]float64{{0.06}, {0.09, 0}, {0.09, 0, 0}, {0.12}}
	attackHitlagFactor    = [][]float64{{0.01}, {0.01, 0}, {0.01, 0, 0}, {0.01}}
	attackOffsets         = []float64{-0.7, -1.0, -0.7, 0.} // y only
	attackHitboxes        = [][]float64{{3., 3.}, {2.0, 3.9}, {2.2, 4.3}, {2.0, 3.9}}
	attackPoiseDMG        = [][]float64{{131.0}, {131.0, 0}, {117.0, 0, 0}, {160.0}}
)

var (
	bikeAttackFrames   [][]int
	bikeAttackHitmarks = []int{31, 32, 29, 35, 41}

	// remove later
	attack     = [][]float64{{1.582}, {0.721, 0.721}, {0.657, 0.657, 0.657}, {2.297}}
	bikeAttack = []float64{1.132, 1.169, 1.383, 1.378, 1.799}
)

const (
	normalHitNum = 4
	bikeHitNum   = 5
)

func init() {
	// Normal attack
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 47)
	attackFrames[0][action.ActionAttack] = 30

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 48)
	attackFrames[1][action.ActionAttack] = 41

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 48)
	attackFrames[2][action.ActionAttack] = 41

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 79)
	attackFrames[3][action.ActionAttack] = 75

	// Skill attack
	bikeAttackFrames = make([][]int, bikeHitNum)

	bikeAttackFrames[0] = frames.InitNormalCancelSlice(bikeAttackHitmarks[0], 53)
	bikeAttackFrames[1] = frames.InitNormalCancelSlice(bikeAttackHitmarks[1], 53)
	bikeAttackFrames[2] = frames.InitNormalCancelSlice(bikeAttackHitmarks[2], 53)
	bikeAttackFrames[3] = frames.InitNormalCancelSlice(bikeAttackHitmarks[3], 53)
	bikeAttackFrames[4] = frames.InitNormalCancelSlice(bikeAttackHitmarks[4], 53)
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() && c.allFireArmamnetsActive {
		return c.bikeAttack(p)
	}
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:          attacks.AttackTagNormal,
		ICDTag:             attacks.ICDTagNormalAttack,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		Element:            attributes.Physical,
		Durability:         25,
		CanBeDefenseHalted: true,
	}
	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult
		ai.PoiseDMG = attackPoiseDMG[c.NormalCounter][i]
		ai.HitlagFactor = attackHitlagFactor[c.NormalCounter][i]
		ai.HitlagHaltFrames = attackHitlagHaltFrame[c.NormalCounter][i] * 60

		ap := combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
			attackHitboxes[c.NormalCounter][1],
		)

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

func (c *char) bikeAttack(_ map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           fmt.Sprintf("Flamestrider Normal Attack %d", c.normalBCounter),
		AttackTag:      attacks.AttackTagNormal,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagNormalAttack,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeBlunt,
		Element:        attributes.Pyro,
		Durability:     25,
		FlatDmg:        c.TotalAtk() * bikeAttack[c.normalBCounter],
		IgnoreInfusion: true,
	}

	ap := combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 2)
	if c.StatusIsActive(crucibleOfDeathAndLifeStatus) {
		ai.FlatDmg += 0.0072 * c.TotalAtk() * float64(c.consumedFightingSpirit)
	}
	ai.FlatDmg += c.c2FlatIncrease(attacks.AttackTagNormal)
	c.Core.QueueAttack(ai, ap, bikeAttackHitmarks[c.normalBCounter], bikeAttackHitmarks[c.normalBCounter])

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAbilFunc(bikeAttackFrames[c.normalBCounter]),
		AnimationLength: bikeAttackFrames[c.normalBCounter][action.InvalidAction],
		CanQueueAfter:   bikeAttackFrames[c.normalBCounter][action.ActionDash],
		State:           action.NormalAttackState,
	}, nil
}
