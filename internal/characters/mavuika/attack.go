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
	attackHitmarks        = [][]int{{21}, {11, 23}, {10, 18, 26}, {28}}
	attackPoiseDMG        = []float64{93.33, 92.72, 115.14, 143.17}
	attackHitlagHaltFrame = []float64{0.09, 0.10, 0.08, .12}
	attackHitboxes        = []float64{2.2, 2.3, 1.8, 3}
	attackOffsets         = []float64{0.5, -1.3, 0.5, -0.8}

	bikeAttackFrames          [][]int
	bikeAttackHitmarks        = []int{21, 22, 27, 14, 41}
	bikeAttackPoiseDMG        = []float64{76.6, 79.1, 93.6, 93.2, 121.7}
	bikeAttackHitlagHaltFrame = []float64{0.09, 0.08, 0.04, 0.03, 0.0}
	bikeAttackHitboxes        = [][]float64{{3.7}, {4}, {3.7}, {5.5, 4.5}, {4.7}}
	bikeAttackOffsets         = []float64{0.5, -1.3, 0.5, -0.8, 1}
)

const normalHitNum = 4
const bikeHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 35) // N1 -> N2
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][1], 44) // N2 -> N3
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][2], 54) // N3 -> N4
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 42) // N4 -> N1

	bikeAttackFrames = make([][]int, bikeHitNum)

	bikeAttackFrames[0] = frames.InitNormalCancelSlice(bikeAttackHitmarks[0], 26) // N1 -> N2
	bikeAttackFrames[1] = frames.InitNormalCancelSlice(bikeAttackHitmarks[1], 35) // N2 -> N3
	bikeAttackFrames[2] = frames.InitNormalCancelSlice(bikeAttackHitmarks[2], 34) // N3 -> N4
	bikeAttackFrames[3] = frames.InitNormalCancelSlice(bikeAttackHitmarks[3], 22) // N4 -> N5
	bikeAttackFrames[4] = frames.InitNormalCancelSlice(bikeAttackHitmarks[4], 63) // N5 -> N1
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.armamentState == bike && c.nightsoulState.HasBlessing() {
		return c.bikeAttack(), nil
	}
	ai := combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:        attacks.AttackTagNormal,
		ICDTag:           attacks.ICDTagNormalAttack,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeBlunt,
		PoiseDMG:         attackPoiseDMG[c.NormalCounter],
		Element:          attributes.Physical,
		Durability:       25,
		Mult:             attack[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor:     0.01,
		HitlagHaltFrames: attackHitlagHaltFrame[c.NormalCounter] * 60,
	}
	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		geometry.Point{Y: attackOffsets[c.NormalCounter]},
		attackHitboxes[c.NormalCounter],
	)

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

// Mavuika's bike normals do not immediately reset after sprinting and will carry the counter upon losing NS Blessing
func (c *char) bikeAttack() action.Info {
	switch c.Core.Player.CurrentState() {
	case action.DashState:
		c.NormalCounter = c.savedNormalCounter
	}
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

	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		geometry.Point{Y: bikeAttackOffsets[c.NormalCounter]},
		bikeAttackHitboxes[c.NormalCounter][0],
	)

	if c.NormalCounter == 3 {
		ap = combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: bikeAttackOffsets[c.NormalCounter]},
			bikeAttackHitboxes[c.NormalCounter][0],
			bikeAttackHitboxes[c.NormalCounter][1],
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

// TODO: charged attack
