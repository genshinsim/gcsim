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
	attackHitmarks        = []int{18, 12, 23}
	attackHitlagHaltFrame = []float64{0.03, 0.03, 0.03}
	attackHitlagFactor    = []float64{0.01, 0.01, 0.01}
	attackDefHalt         = []bool{true, true, true}
	attackHitboxes        = [][]float64{{3., 3.}, {3.9, 2.}, {4.3, 2.2}}
)

var (
	skillAttackFrames        [][]int
	skillAttackHitmarks      = [][]int{{10, 20}, {10, 19}}
	skillAttackAngularTravel = []float64{77.709, 63.436}
)

const (
	normalHitNum = 3
	skillHitNum  = 2

	loopShotNSGenDelay = 1

	angularVelocity   = 1.586  // degrees per frame
	blindSpotBoundary = 33.387 // +- degrees from the "center" of the blind spot
)

func init() {
	// Normal attack
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 24) // N1 -> CA
	attackFrames[0][action.ActionAttack] = 19

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 27)
	attackFrames[1][action.ActionAttack] = 12

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 42)

	// Skill attack
	skillAttackFrames = make([][]int, skillHitNum)

	skillAttackFrames[0] = frames.InitNormalCancelSlice(skillAttackHitmarks[0][1], 18)
	skillAttackFrames[0][action.ActionSkill] = 11
	skillAttackFrames[0][action.ActionBurst] = 10

	skillAttackFrames[1] = frames.InitNormalCancelSlice(skillAttackHitmarks[1][1], 17)
	skillAttackFrames[1][action.ActionSkill] = 10
	skillAttackFrames[1][action.ActionBurst] = 10
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

	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 0.5)

	c.Core.Tasks.Add(c.loopShotGenerateNSPoints, loopShotNSGenDelay)

	c.Core.QueueAttack(ai, ap, skillAttackHitmarks[c.normalSCounter][0]+travel, skillAttackHitmarks[c.normalSCounter][0]+travel, c.desolationCB, c.c2ResShredCB)
	c.Core.QueueAttack(ai, ap, skillAttackHitmarks[c.normalSCounter][1]+travel, skillAttackHitmarks[c.normalSCounter][1]+travel, c.desolationCB, c.c2ResShredCB)

	defer c.AdvanceNormalIndex()

	direction, ok := p["direction"]
	fmt.Println(c.Core.F, direction, ok)
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
			fmt.Println("\n", c.Core.F, "killed blind spot", int(math.Abs(boundary-c.characterAngularPosition)/angularVelocity), "frames after attack start")
		}, int(math.Abs(boundary-c.characterAngularPosition)/angularVelocity))
	}
	c.QueueCharTask(func() {
		c.characterAngularPosition = NormalizeAngle(c.characterAngularPosition + float64(direction)*skillAttackAngularTravel[c.normalSCounter])
	}, skillAttackFrames[c.normalSCounter][action.ActionBurst]) // CHANGE this when frames are ready

	return action.Info{
		Frames:          frames.NewAbilFunc(skillAttackFrames[c.normalSCounter]),
		AnimationLength: skillAttackFrames[c.normalSCounter][action.InvalidAction],
		CanQueueAfter:   skillAttackFrames[c.normalSCounter][action.ActionBurst],
		State:           action.SkillState,
	}, nil
}

func (c *char) NextMoveIsInBlindSpot(direction int) (bool, float64) {
	fmt.Println("\n", c.Core.F)
	if c.blindSpotAngularPosition == -1 {
		fmt.Println("Kinich attacks, but blind spot is absent")
		return false, -1
	}
	// Calculate the sector boundaries and normalize them
	lowerBoundary := NormalizeAngle(c.blindSpotAngularPosition - blindSpotBoundary)
	upperBoundary := NormalizeAngle(c.blindSpotAngularPosition + blindSpotBoundary)

	targetPos := NormalizeAngle(c.characterAngularPosition + float64(direction)*skillAttackAngularTravel[c.normalSCounter])
	fmt.Println("Kinich attacks and the blind spot is present. His position:", c.characterAngularPosition, ", blind spot position:", c.blindSpotAngularPosition)
	fmt.Println("And the blind spot bounds are:", lowerBoundary, upperBoundary)
	fmt.Println("Kinich is going to move by", float64(direction)*skillAttackAngularTravel[c.normalSCounter], "so his position will be:", targetPos)

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
	fmt.Println("Target in sector:", targetInSector)

	// Determine which boundary is first crossed based on direction
	if targetInSector {
		if direction == -1 {
			return true, upperBoundary // Crosses upper boundary moving clockwise
		}
		return true, lowerBoundary // Crosses lower boundary moving counterclockwise
	}
	lowerBoundaryStart := NormalizeAngle(lowerBoundary - c.characterAngularPosition)
	targetUpperBoundary := NormalizeAngle(targetPos - upperBoundary)
	startUpperBoundary := NormalizeAngle(c.characterAngularPosition - upperBoundary)
	lowerBoundaryTarget := NormalizeAngle(lowerBoundary - targetPos)

	if lowerBoundaryStart > 0 && targetUpperBoundary > 0 && lowerBoundaryStart+targetUpperBoundary+blindSpotBoundary*2-skillAttackAngularTravel[c.normalSCounter] < 0.0001 {
		fmt.Println("Kinich jumped over the blind spot moving counter clockwise")
		return true, lowerBoundary
	}
	if startUpperBoundary > 0 && lowerBoundaryTarget > 0 && startUpperBoundary+lowerBoundaryTarget+blindSpotBoundary*2-skillAttackAngularTravel[c.normalSCounter] < 0.0001 {
		fmt.Println("Kinich jumped over the blind spot moving clockwise")
		return true, upperBoundary
	}

	return false, -1 // Point does not enter the sector after moving
}

// Normalize an angle to be within [0, 360)
func NormalizeAngle(angle float64) float64 {
	for angle < 0 {
		angle += 360
	}
	for angle >= 360 {
		angle -= 360
	}
	return angle
}
