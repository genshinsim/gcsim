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
	attackFrames          [][]int
	attackHitmarks        = [][]int{{21}, {14, 26}, {28, 33, 39}, {30}}
	attackPoiseDMG        = []float64{107.0, 48.8, 44.43, 155.4}
	attackHitlagHaltFrame = []float64{0.09, 0.05, 0.0, 0.1}
	attackHitboxes        = [][]float64{{2.2}, {3.3, 4.3}, {2.8, 5.0}, {3.2}}
	attackOffsets         = []float64{0.5, -1.3, 0.5, -0.8}

	bikeAttackFrames          [][]int
	bikeAttackHitmarks        = []int{19, 24, 31, 13, 37}
	bikeAttackPoiseDMG        = []float64{76.6, 79.1, 93.6, 93.2, 121.7}
	bikeAttackHitlagHaltFrame = []float64{0.09, 0.08, 0.04, 0.03, 0.0}
	bikeAttackHitboxes        = [][]float64{{3.2}, {3.5}, {3.0}, {3.5, 4.5}, {4.0}}
	bikeAttackBurstHitboxes   = [][]float64{{3.7}, {4.0}, {3.7}, {4.5, 5.5}, {4.7}}
	bikeAttackOffsets         = []float64{1.5, 1.85, 0.5, -0.8, 1}
)

const normalHitNum = 4
const bikeHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 40) // N1 -> Walk
	attackFrames[0][action.ActionAttack] = 31                                // N1 -> N2
	attackFrames[0][action.ActionCharge] = 31                                // N1 -> CA

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][1], 50) // N2 -> Walk
	attackFrames[1][action.ActionAttack] = 42                                // N2 -> N3
	attackFrames[1][action.ActionCharge] = 42                                // N2 -> CA

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][2], 59) // N3 -> Walk
	attackFrames[2][action.ActionAttack] = 46                                // N3 -> N4
	attackFrames[2][action.ActionCharge] = 47                                // N3 -> CA

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 61) // N4 -> Walk
	attackFrames[3][action.ActionAttack] = 60                                // N4 -> N1
	attackFrames[3][action.ActionCharge] = 61                                // N4 -> CA

	bikeAttackFrames = make([][]int, bikeHitNum)

	bikeAttackFrames[0] = frames.InitNormalCancelSlice(bikeAttackHitmarks[0], 39) // N1 -> Walk
	bikeAttackFrames[0][action.ActionAttack] = 23                                 // N1 -> N2
	bikeAttackFrames[0][action.ActionCharge] = 33                                 // N1 -> CA

	bikeAttackFrames[1] = frames.InitNormalCancelSlice(bikeAttackHitmarks[1], 46) // N2 -> Walk
	bikeAttackFrames[1][action.ActionAttack] = 32                                 // N2 -> N3
	bikeAttackFrames[1][action.ActionCharge] = 35                                 // N2 -> CA

	bikeAttackFrames[2] = frames.InitNormalCancelSlice(bikeAttackHitmarks[2], 50) // N3 -> Walk
	bikeAttackFrames[2][action.ActionAttack] = 35                                 // N3 -> N4
	bikeAttackFrames[2][action.ActionCharge] = 40                                 // N3 -> CA

	bikeAttackFrames[3] = frames.InitNormalCancelSlice(bikeAttackHitmarks[3], 44) // N4 -> Walk
	bikeAttackFrames[3][action.ActionAttack] = 27                                 // N4 -> N5
	bikeAttackFrames[3][action.ActionCharge] = 29                                 // N4 -> CA

	bikeAttackFrames[4] = frames.InitNormalCancelSlice(bikeAttackHitmarks[4], 70) // N5 -> Walk
	bikeAttackFrames[4][action.ActionAttack] = 68                                 // N5 -> N1
	bikeAttackFrames[4][action.ActionCharge] = 64                                 // N5 -> CA
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.armamentState == bike && c.nightsoulState.HasBlessing() {
		return c.bikeAttack(), nil
	}
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:          attacks.AttackTagNormal,
		ICDTag:             attacks.ICDTagNormalAttack,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		PoiseDMG:           attackPoiseDMG[c.NormalCounter],
		Element:            attributes.Physical,
		Durability:         25,
		Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter] * 60,
		CanBeDefenseHalted: true,
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
	default: // circle
		ap = combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
		)
	}

	for _, delay := range attackHitmarks[c.NormalCounter] {
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

// Mavuika's bike normals do not immediately reset after sprinting and will carry the counter upon losing NS Blessing (Handled in Dash)
func (c *char) bikeAttack() action.Info {
	delay := bikeAttackHitmarks[c.NormalCounter]
	ai := combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             fmt.Sprintf("Flamestrider Normal %v", c.NormalCounter),
		AttackTag:        attacks.AttackTagNormal,
		AdditionalTags:   []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:           attacks.ICDTagMavuikaFlamestrider,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeBlunt,
		PoiseDMG:         bikeAttackPoiseDMG[c.NormalCounter],
		Element:          attributes.Pyro,
		Durability:       25,
		Mult:             skillAttack[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor:     0.01,
		HitlagHaltFrames: bikeAttackHitlagHaltFrame[c.NormalCounter] * 60,
		IgnoreInfusion:   true,
	}

	hitboxes := bikeAttackHitboxes
	if c.StatusIsActive(burstKey) {
		hitboxes = bikeAttackBurstHitboxes
	}

	var ap combat.AttackPattern
	switch {
	case len(hitboxes[c.NormalCounter]) == 2: // box
		ap = combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: bikeAttackOffsets[c.NormalCounter]},
			hitboxes[c.NormalCounter][0],
			hitboxes[c.NormalCounter][1],
		)
	default: // circle
		ap = combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: bikeAttackOffsets[c.NormalCounter]},
			hitboxes[c.NormalCounter][0],
		)
	}

	c.QueueCharTask(func() {
		ai.FlatDmg = c.burstBuffNA() + c.c2BikeNA()
		c.Core.QueueAttack(ai, ap, 0, 0)
		c.reduceNightsoulPoints(1)
	}, delay)

	defer func() {
		c.AdvanceNormalIndex()
		c.savedNormalCounter = c.NormalCounter
	}()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, bikeAttackFrames),
		AnimationLength: bikeAttackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   bikeAttackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}
