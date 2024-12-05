package kinich

import (
	"fmt"
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var (
	attackFrames          [][]int
	attackHitmarks        = []int{21, 22, 44}
	attackHitlagHaltFrame = []float64{0.06, 0.09, 0.12}
	attackHitlagFactor    = []float64{0.01, 0.01, 0.01}
	attackDefHalt         = []bool{true, true, true}
	attackOffsets         = []float64{-0.7, -1.0, 0.0} // y only
	attackHitboxes        = [][]float64{{3., 3.}, {2.0, 3.9}, {2.2, 4.3}}
	attackPoiseDMG        = []float64{131.0, 117.0, 160.0}
)

var (
	skillAttackFrames        [][]int
	skillAttackHitmarks      = [][]int{{30, 38}, {31, 38}}
	skillAttackAngularTravel = []float64{70., 70.}
)

const (
	normalHitNum = 3
	skillHitNum  = 2

	angularVelocity   = 70. / 40 // degrees per frame
	blindSpotBoundary = 35.      // +- degrees from the "center" of the blind spot
)

func init() {
	// Normal attack
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 47) // N1 -> Walk
	attackFrames[0][action.ActionAttack] = 30

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 48) // N2 -> Walk
	attackFrames[1][action.ActionAttack] = 41

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 79) // N3 -> Walk
	attackFrames[2][action.ActionAttack] = 75

	// Skill attack
	skillAttackFrames = make([][]int, skillHitNum)

	skillAttackFrames[0] = frames.InitNormalCancelSlice(skillAttackHitmarks[0][1], 53) // N1(E) -> Swap
	skillAttackFrames[0][action.ActionAttack] = 40
	skillAttackFrames[0][action.ActionSkill] = 33
	skillAttackFrames[0][action.ActionBurst] = 35
	skillAttackFrames[0][action.ActionDash] = 26
	skillAttackFrames[0][action.ActionJump] = 32
	skillAttackFrames[0][action.ActionWalk] = 50

	skillAttackFrames[1] = frames.InitNormalCancelSlice(skillAttackHitmarks[1][1], 53) // N2(E) -> Swap
	skillAttackFrames[1][action.ActionAttack] = 40
	skillAttackFrames[1][action.ActionSkill] = 33
	skillAttackFrames[1][action.ActionBurst] = 32
	skillAttackFrames[1][action.ActionDash] = 21
	skillAttackFrames[1][action.ActionJump] = 34
	skillAttackFrames[1][action.ActionWalk] = 51
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		return c.skillAttack(p)
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
		HitlagFactor:       attackHitlagFactor[c.NormalCounter],
		HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter] * 60,
		CanBeDefenseHalted: attackDefHalt[c.NormalCounter],
	}

	ap := combat.NewBoxHitOnTarget(
		c.Core.Combat.Player(),
		geometry.Point{Y: attackOffsets[c.NormalCounter]},
		attackHitboxes[c.NormalCounter][0],
		attackHitboxes[c.NormalCounter][1],
	)

	c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter], attackHitmarks[c.NormalCounter])

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}, nil
}

func (c *char) skillAttack(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           fmt.Sprintf("Loop Shot %d", c.normalSCounter),
		AttackTag:      attacks.AttackTagElementalArt,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagKinichLoopShot,
		ICDGroup:       attacks.ICDGroupKinichLoopShot,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Dendro,
		Durability:     25,
		Mult:           loopShot[c.TalentLvlSkill()],
		IgnoreInfusion: true,
	}

	c.loopShotGenerateNSPoints()

	ap := combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 0.5)
	c.Core.QueueAttack(ai, ap, skillAttackHitmarks[c.normalSCounter][0], skillAttackHitmarks[c.normalSCounter][0], c.a1CB, c.c2ResShredCB)
	c.Core.QueueAttack(ai, ap, skillAttackHitmarks[c.normalSCounter][1], skillAttackHitmarks[c.normalSCounter][1], c.a1CB, c.c2ResShredCB)
	c.Core.Tasks.Add(c.c4, skillAttackHitmarks[c.normalSCounter][0])

	defer c.AdvanceNormalIndex()

	direction, ok := p["direction"]
	if !ok {
		direction = -1
	}
	switch direction {
	case 1:
	case -1:
	default:
		return action.Info{}, fmt.Errorf("%v, %v: Wrong value of direction: %v, should be 1 or -1", c.Core.F, c.Base.Key, direction)
	}
	cross, boundary := c.NextMoveIsInBlindSpot(direction)
	if cross {
		time := math.Abs(NormalizeAngle180(boundary-c.characterAngularPosition)) / angularVelocity
		c.QueueCharTask(func() {
			c.nightsoulState.GeneratePoints(4)
			c.blindSpotAngularPosition = -1
		}, int(time))
	}
	c.QueueCharTask(func() {
		c.characterAngularPosition = NormalizeAngle360(c.characterAngularPosition + float64(direction)*skillAttackAngularTravel[c.normalSCounter])
	}, skillAttackFrames[c.normalSCounter][action.ActionDash]) // Update on earliest possible cancel

	return action.Info{
		Frames:          frames.NewAbilFunc(skillAttackFrames[c.normalSCounter]),
		AnimationLength: skillAttackFrames[c.normalSCounter][action.InvalidAction],
		CanQueueAfter:   skillAttackFrames[c.normalSCounter][action.ActionDash],
		State:           action.NormalAttackState,
	}, nil
}

func (c *char) NextMoveIsInBlindSpot(direction int) (bool, float64) {
	if c.blindSpotAngularPosition == -1 {
		return false, -1
	}

	// Calculate the sector boundaries and normalize them
	lowerBoundary := NormalizeAngle360(c.characterAngularPosition)
	upperBoundary := NormalizeAngle360(c.characterAngularPosition + float64(direction)*skillAttackAngularTravel[c.normalSCounter])

	lowerBlindBoundary := NormalizeAngle360(c.blindSpotAngularPosition - blindSpotBoundary)
	upperBlindBoundary := NormalizeAngle360(c.blindSpotAngularPosition + blindSpotBoundary)

	if direction == -1 {
		lowerBoundary, upperBoundary = upperBoundary, lowerBoundary
		lowerBlindBoundary, upperBlindBoundary = upperBlindBoundary, lowerBlindBoundary
	}

	// Helper function to check if an angle is within the circular sector
	isInSector := func(angle float64) bool {
		if lowerBoundary < upperBoundary {
			return angle >= lowerBoundary && angle <= upperBoundary
		}
		// Handles wrap-around sector cases (e.g., sector from 350 to 10 degrees)
		return angle >= lowerBoundary || angle <= upperBoundary
	}

	if isInSector(lowerBlindBoundary) {
		return true, lowerBlindBoundary
	}
	if isInSector(upperBlindBoundary) {
		return true, upperBlindBoundary
	}
	return false, -1
}
