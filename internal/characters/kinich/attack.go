package kinich

import (
	"fmt"
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	attackFrames          [][]int
	attackHitmarks        = []int{21, 22, 44}
	attackHitlagHaltFrame = []float64{0.03, 0.03, 0.03}
	attackHitlagFactor    = []float64{0.01, 0.01, 0.01}
	attackDefHalt         = []bool{true, true, true}
	attackHitboxes        = [][]float64{{3., 3.}, {3.9, 2.}, {4.3, 2.2}}
)

var (
	skillAttackFrames        [][]int
	skillAttackHitmarks      = [][]int{{30, 38}, {31, 38}}
	skillAttackAngularTravel = []float64{70., 70.}
)

const (
	normalHitNum = 3
	skillHitNum  = 2

	loopShotNSGenDelay = 0

	angularVelocity   = 70. / 40 // degrees per frame
	blindSpotBoundary = 35.      // +- degrees from the "center" of the blind spot
)

func init() {
	// Normal attack
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 47)
	attackFrames[0][action.ActionAttack] = 30

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 48)
	attackFrames[1][action.ActionAttack] = 41

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 79)
	attackFrames[2][action.ActionAttack] = 75

	// Skill attack
	skillAttackFrames = make([][]int, skillHitNum)

	skillAttackFrames[0] = frames.InitNormalCancelSlice(skillAttackHitmarks[0][1], 53)
	skillAttackFrames[0][action.ActionSkill] = 33
	skillAttackFrames[0][action.ActionBurst] = 33
	skillAttackFrames[0][action.ActionDash] = 24
	skillAttackFrames[0][action.ActionJump] = 32
	skillAttackFrames[0][action.ActionWalk] = 51

	skillAttackFrames[1] = frames.InitNormalCancelSlice(skillAttackHitmarks[1][1], 53)
	skillAttackFrames[1][action.ActionSkill] = 33
	skillAttackFrames[1][action.ActionBurst] = 33
	skillAttackFrames[1][action.ActionDash] = 24
	skillAttackFrames[1][action.ActionJump] = 32
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
		StrikeType:         attacks.StrikeTypeSlash,
		Element:            attributes.Physical,
		Durability:         25,
		Mult:               attack[c.normalSCounter][c.TalentLvlAttack()],
		HitlagFactor:       attackHitlagFactor[c.NormalCounter],
		HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter] * 60,
		CanBeDefenseHalted: attackDefHalt[c.NormalCounter],
	}

	ap := combat.NewBoxHitOnTarget(
		c.Core.Combat.Player(),
		nil,
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
		ICDTag:         attacks.ICDTagKinichLoopShot,
		ICDGroup:       attacks.ICDGroupKinichLoopShot,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Dendro,
		Durability:     25,
		Mult:           loopShot[c.TalentLvlSkill()],
		IgnoreInfusion: true,
	}

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 0.5)

	c.Core.Tasks.Add(c.loopShotGenerateNSPoints, loopShotNSGenDelay)

	c.Core.QueueAttack(ai, ap, skillAttackHitmarks[c.normalSCounter][0], skillAttackHitmarks[c.normalSCounter][0], c.desolationCB, c.c2ResShredCB)
	c.Core.QueueAttack(ai, ap, skillAttackHitmarks[c.normalSCounter][1], skillAttackHitmarks[c.normalSCounter][1], c.desolationCB, c.c2ResShredCB)

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
		c.QueueCharTask(func() {
			c.nightsoulState.GeneratePoints(4)
			c.blindSpotAngularPosition = -1
		}, int(math.Abs(NormalizeAngle180(boundary-c.characterAngularPosition))/angularVelocity))
	}
	c.QueueCharTask(func() {
		c.characterAngularPosition = NormalizeAngle360(c.characterAngularPosition + float64(direction)*skillAttackAngularTravel[c.normalSCounter])
	}, skillAttackFrames[c.normalSCounter][action.ActionDash]) // Update on earliest possible cancel

	return action.Info{
		Frames:          frames.NewAbilFunc(skillAttackFrames[c.normalSCounter]),
		AnimationLength: skillAttackFrames[c.normalSCounter][action.InvalidAction],
		CanQueueAfter:   skillAttackFrames[c.normalSCounter][action.ActionBurst],
		State:           action.SkillState,
	}, nil
}

func (c *char) NextMoveIsInBlindSpot(direction int) (bool, float64) {
	if c.blindSpotAngularPosition == -1 {
		return false, -1
	}
	// Calculate the sector boundaries and normalize them
	lowerBoundary := NormalizeAngle360(c.blindSpotAngularPosition - blindSpotBoundary)
	upperBoundary := NormalizeAngle360(c.blindSpotAngularPosition + blindSpotBoundary)

	targetPos := NormalizeAngle360(c.characterAngularPosition + float64(direction)*skillAttackAngularTravel[c.normalSCounter])

	// Helper function to check if an angle is within the circular sector
	isInSector := func(angle float64) bool {
		if lowerBoundary < upperBoundary {
			return angle >= lowerBoundary && angle <= upperBoundary
		}
		// Handles wrap-around sector cases (e.g., sector from 350 to 10 degrees)
		return angle >= lowerBoundary || angle <= upperBoundary
	}

	// Determine if the target position is within the sector
	targetInSector := isInSector(targetPos)

	// Determine which boundary is first crossed based on direction
	if targetInSector {
		if direction == -1 {
			return true, upperBoundary // Crosses upper boundary moving clockwise
		}
		return true, lowerBoundary // Crosses lower boundary moving counterclockwise
	}
	lowerBoundaryStart := NormalizeAngle360(lowerBoundary - c.characterAngularPosition)
	targetUpperBoundary := NormalizeAngle360(targetPos - upperBoundary)
	startUpperBoundary := NormalizeAngle360(c.characterAngularPosition - upperBoundary)
	lowerBoundaryTarget := NormalizeAngle360(lowerBoundary - targetPos)

	if lowerBoundaryStart > 0 && targetUpperBoundary > 0 && lowerBoundaryStart+targetUpperBoundary+blindSpotBoundary*2-skillAttackAngularTravel[c.normalSCounter] < 0.0001 {
		return true, lowerBoundary
	}
	if startUpperBoundary > 0 && lowerBoundaryTarget > 0 && startUpperBoundary+lowerBoundaryTarget+blindSpotBoundary*2-skillAttackAngularTravel[c.normalSCounter] < 0.0001 {
		return true, upperBoundary
	}

	return false, -1 // Point does not enter the sector after moving
}
