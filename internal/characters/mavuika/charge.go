package mavuika

import (
	"fmt"
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var chargeFrames []int
var bikeChargeFrames []int
var bikeChargeFinalFrames []int
var bikeChargeAttackHittableList []HittableEntity

// Bike CA has 14f of c0 before going into CA

// Minimum CA time before CAF anim is 50f
var bikeChargeAttackMinimumDuration = 50
var bikeChargeAttackStartupHitmark = 35

// Maximum CA time before CAF anim is 375f
var bikeChargeAttackMaximumDuration = 375
var bikeChargeFinalHitmark = 45

// TODO: Replicate frames 35-46 of the CA more accurately
// var bikeSpinInitialFrames = 11
// var bikeSpinInitialAngularVelocity = float64(-90 / 11)
// spin velocity varies by current angle
var bikeSpinQuadrantAngularVelocity = []float64{-90 / 9, -90 / 7, -90 / 15, -90 / 14} // Quadrant 4, 3, 2, 1
var bikeSpinQuadrantFrames = []int{9, 7, 15, 14}                                      // Quadrant 4, 3, 2, 1

const chargeHitmark = 40
const bikeChargeAttackICD = 42         // Minimum time between CA hits
const bikeChargeAttackSpinFrames = 45  // One revolution every 45f
const bikeChargeAttackHitboxRadius = 3 // Placeholder
const bikeChargeAttackSpinOffset = 4.0 // Estimated center of hitbox from Mav origin

func init() {
	chargeFrames = frames.InitAbilSlice(48)
	chargeFrames[action.ActionBurst] = 50
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionSwap] = 50
	chargeFrames[action.ActionWalk] = 60

	bikeChargeFrames = frames.InitAbilSlice(bikeChargeAttackStartupHitmark)
	bikeChargeFrames[action.ActionCharge] = 0
	bikeChargeFrames[action.ActionBurst] = bikeChargeAttackStartupHitmark
	bikeChargeFrames[action.ActionDash] = bikeChargeAttackStartupHitmark
	bikeChargeFrames[action.ActionJump] = bikeChargeAttackStartupHitmark
	bikeChargeFrames[action.ActionSwap] = bikeChargeAttackStartupHitmark

	bikeChargeFinalFrames = frames.InitAbilSlice(bikeChargeFinalHitmark)
	bikeChargeFinalFrames[action.ActionBurst] = bikeChargeFinalHitmark
	bikeChargeFinalFrames[action.ActionDash] = bikeChargeFinalHitmark
	bikeChargeFinalFrames[action.ActionJump] = bikeChargeFinalHitmark
	bikeChargeFinalFrames[action.ActionSwap] = bikeChargeFinalHitmark
}

// Charge state struct
type ChargeState struct {
	StartFrame int
	CAtkFrames int
	LastHitF   int
}

type HittableEntity struct {
	Entity     combat.Target
	isOneTick  bool   // Does entity get destroyed after a single maxHitCount?
	CollFrames [2]int // Frames of the CA spin on which collision happens
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if c.armamentState == bike && c.nightsoulState.HasBlessing() {
		return c.bikeCharge(p)
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNormalAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   120.0,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: -1.8},
			2,
			4.5,
		),
		chargeHitmark,
		chargeHitmark,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmark,
		State:           action.ChargeAttackState,
	}, nil
}

// This starts the CA, then goes to a loop handler for duration calc
func (c *char) bikeCharge(p map[string]int) (action.Info, error) {
	c.BuildBikeChargeAttackHittableTargetList()

	// Check if a continuing CA or new
	if c.Core.Player.CurrentState() != action.ChargeAttackState {
		c.caState = ChargeState{}
		c.caState.StartFrame = c.Core.F
		c.caState.LastHitF = 0
	}
	// Add any existing CA frames
	c.caState.CAtkFrames = (c.Core.F - c.caState.StartFrame)

	// Parameters for tuning CA
	durationCA := p["hold"]
	chargeCount := p["hits"]
	final := p["final"]

	// TODO: Allow more ways to get into the final attack?
	if final == 1 {
		return c.bikeChargeFinalAttack(), nil
	}

	switch {
	case durationCA > 0:
		// c.caState.CAtkFrames += durationCA
		// Default max hold time is 6.25s/375f
		if durationCA > bikeChargeAttackMaximumDuration {
			c.caState.CAtkFrames = durationCA - bikeChargeAttackMaximumDuration // Cap additional hold time to maximum
			durationCA = bikeChargeAttackMaximumDuration
		} else if (durationCA + c.caState.StartFrame) < bikeChargeAttackStartupHitmark {
			durationCA = bikeChargeAttackStartupHitmark - c.caState.StartFrame
			c.caState.CAtkFrames = durationCA
		}
		// Hold CA logic
		c.Core.Log.NewEventBuildMsg(glog.LogSimEvent, c.Index,
			fmt.Sprintf("Last hit is %d | Start frame is %d | Duration is %d | State frames is %d",
				c.caState.LastHitF, c.caState.StartFrame, durationCA, c.caState.CAtkFrames))
		c.HoldBikeChargeAttack(durationCA)
	case chargeCount > 0:
		// CA count logic
		// No idea why this error check always returns true as if the target list is empty
		// isTargetForCountsHittable := false
		// if len(bikeChargeAttackHittableList) == 0 {
		// 	return action.Info{}, errors.New("no valid targets within flamestrider area")
		// }
		// for _, t := range bikeChargeAttackHittableList {
		// 	if t.Entity == c.Core.Combat.PrimaryTarget() {
		// 		isTargetForCountsHittable = true
		// 		break
		// 	}
		// }
		// if !isTargetForCountsHittable {
		// 	return action.Info{}, errors.New("primary target is not within flamestrider area")
		// }
		durationCA = c.CountBikeChargeAttack(chargeCount)
	case durationCA == 0 && final == 0 && chargeCount == 0:
		// Default to single CA if nothing specified
		durationCA = c.CountBikeChargeAttack(1)
	}

	c.Core.Tasks.Add(func() {
		c.bikeChargeAttackHook()
	}, bikeChargeAttackStartupHitmark-1)

	c.Core.Tasks.Add(func() {
		c.bikeChargeAttackUnhook()
	}, durationCA)

	return action.Info{
		Frames:          func(next action.Action) int { return durationCA },
		AnimationLength: durationCA,
		CanQueueAfter:   durationCA,
		State:           action.ChargeAttackState,
	}, nil
}

// For given CA length, calculate hits on each target in hittable list
func (c *char) HoldBikeChargeAttack(CAtkFrames int) {
	lastPrimaryHitF := 0
	for i := 0; i < len(bikeChargeAttackHittableList); i++ {
		t := bikeChargeAttackHittableList[i]
		// First 11f of CA are a bit inaccurate, should maxHitCount further left
		hitFrames := c.CalculateValidCollisionFrames(CAtkFrames, t.CollFrames)

		if len(hitFrames) > 0 {
			for _, f := range hitFrames {
				c.Core.Tasks.Add(func() {
					ai := c.GetBikeChargeAttackAttackInfo()
					c.Core.QueueAttack(ai, combat.NewSingleTargetHit(t.Entity.Key()), 0, 0)
				}, f)
				c.Core.Log.NewEventBuildMsg(glog.LogSimEvent, c.Index, fmt.Sprintf("Valid hit frame: %d", f))
				if t.Entity == c.Core.Combat.PrimaryTarget() {
					lastPrimaryHitF = f + (c.caState.CAtkFrames - c.caState.LastHitF)
				}
			}
		}
		// c.Core.Log.NewEventBuildMsg(glog.LogSimEvent, c.Index, fmt.Sprintf("Target %d Spin Collision Frames: %d - %d", i, t.CollFrames[0], t.CollFrames[1]))
	}
	c.caState.LastHitF += lastPrimaryHitF
}

// For given maxHitCount count, calculate maxHitCount timings on targets and return CA duration
func (c *char) CountBikeChargeAttack(maxHitCount int) int {
	dur := bikeChargeAttackMaximumDuration
	hitCounter := 0

	for i := 0; i < len(bikeChargeAttackHittableList); i++ {
		t := bikeChargeAttackHittableList[i]
		if t.Entity != c.Core.Combat.PrimaryTarget() {
			continue
		}

		// First 11f of CA are a bit inaccurate, should maxHitCount further left
		hitFrames := c.CalculateValidCollisionFrames(dur, t.CollFrames)

		if len(hitFrames) > 0 {
			for _, f := range hitFrames {
				hitCounter++
				if hitCounter >= maxHitCount {
					dur = f
					break
				}
			}
		}
		if hitCounter >= maxHitCount {
			break
		}
	}

	for i := 0; i < len(bikeChargeAttackHittableList); i++ {
		t := bikeChargeAttackHittableList[i]

		// First 11f of CA are a bit inaccurate, should maxHitCount further left
		hitFrames := c.CalculateValidCollisionFrames(dur, t.CollFrames)

		if len(hitFrames) > 0 {
			for _, f := range hitFrames {
				c.Core.Tasks.Add(func() {
					ai := c.GetBikeChargeAttackAttackInfo()
					c.Core.QueueAttack(ai, combat.NewSingleTargetHit(t.Entity.Key()), 0, 0)
				}, f)
				c.Core.Log.NewEventBuildMsg(glog.LogSimEvent, c.Index, fmt.Sprintf("Valid hit frame: %d", f))
			}
		}
		c.Core.Log.NewEventBuildMsg(glog.LogSimEvent, c.Index, fmt.Sprintf("Target %d Spin Collision Frames: %d - %d", i, t.CollFrames[0], t.CollFrames[1]))
	}
	// Used when the CA started between hits (hold/hit method mix)
	spinFramesOffset := c.caState.CAtkFrames - c.caState.LastHitF
	c.caState.LastHitF += dur + spinFramesOffset
	return dur
}

func (c *char) bikeChargeFinalAttack() action.Info {
	var adjustedBikeChargeFinalHitmark = bikeChargeFinalHitmark
	bikeChargeAttackElapsedTime := c.caState.StartFrame - c.Core.F
	if bikeChargeAttackElapsedTime < bikeChargeAttackMinimumDuration {
		adjustedBikeChargeFinalHitmark += (50 - bikeChargeAttackElapsedTime)
	}

	c.HoldBikeChargeAttack(50 - bikeChargeAttackElapsedTime)

	c.QueueCharTask(func() {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               "Flamestrider Charged Attack (Final)",
			AttackTag:          attacks.AttackTagExtra,
			AdditionalTags:     []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
			ICDTag:             attacks.ICDTagMavuikaFlamestrider,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeBlunt,
			PoiseDMG:           12.0,
			Element:            attributes.Pyro,
			HitlagFactor:       0.01,
			HitlagHaltFrames:   0.04 * 60,
			CanBeDefenseHalted: true,
			Durability:         25,
			Mult:               skillChargeFinal[c.TalentLvlSkill()],
			IgnoreInfusion:     true,
			FlatDmg:            c.burstBuffCA() + c.c2BikeCA(),
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: 1},
				4,
			),
			0,
			0,
		)
	}, adjustedBikeChargeFinalHitmark)

	return action.Info{
		Frames:          frames.NewAbilFunc(bikeChargeFinalFrames),
		AnimationLength: adjustedBikeChargeFinalHitmark,
		CanQueueAfter:   adjustedBikeChargeFinalHitmark,
		State:           action.ChargeAttackState,
	}
}

func (c *char) GetBikeChargeAttackAttackInfo() combat.AttackInfo {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Flamestrider Charged Attack (Cyclic)",
		AttackTag:      attacks.AttackTagExtra,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagMavuikaFlamestrider,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeBlunt,
		PoiseDMG:       60.0,
		Element:        attributes.Pyro,
		// HitlagFactor:     0.01,
		// HitlagHaltFrames: 0.03 * 60,
		CanBeDefenseHalted: false,
		Durability:         25,
		Mult:               skillCharge[c.TalentLvlSkill()],
		IgnoreInfusion:     true,
		FlatDmg:            c.burstBuffCA() + c.c2BikeCA(),
	}
	return ai
}

// Not sure on the scope of this yet, should definitely trigger on nightsoul stuff
// func (c *char) exitBikeChargeAttack() {

// 	c.bikeChargeAttackUnhook()
// 	// Use switch eventually to determine if CAF should trigger
// 	if c.Core.Player.CurrentState() == action.Idle {
// 	}
// }

func (c *char) BuildBikeChargeAttackHittableTargetList() {
	bikeChargeAttackHittableList = bikeChargeAttackHittableList[:0]
	c.buildValidTargetList()
	c.buildValidGadgetList()
}

func (c *char) buildValidTargetList() {
	c.Core.Combat.Player().Pos()
	enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 8), nil)
	for _, v := range enemies {
		if v == nil {
			continue
		}
		// Calculate start and ending frames for collision
		collisionFrames := [2]int{-1, -1}
		isIntersecting := c.BikeHitboxIntersectionAngles(v, collisionFrames[:])

		if isIntersecting {
			bikeChargeAttackHittableList = append(bikeChargeAttackHittableList, HittableEntity{
				Entity:     combat.Target(v),
				isOneTick:  false,
				CollFrames: collisionFrames,
			})
		}
	}
	// c.Core.Log.NewEventBuildMsg(glog.LogSimEvent, c.Index,
	// 	fmt.Sprintf("There are %v valid targets", len(bikeChargeAttackHittableList)))
}

// Gadgets are gonna be problematic
func (c *char) buildValidGadgetList() {
	gadgets := c.Core.Combat.GadgetsWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 8), nil)
	for _, v := range gadgets {
		if v == nil {
			continue
		}
		if v.GadgetTyp() == combat.GadgetTypDendroCore {
			// Calculate start and ending frames for collision
			collisionFrames := [2]int{-1, -1}
			isIntersecting := c.BikeHitboxIntersectionAngles(v, collisionFrames[:])

			if isIntersecting {
				bikeChargeAttackHittableList = append(bikeChargeAttackHittableList, HittableEntity{
					Entity:     combat.Target(v),
					isOneTick:  true,
					CollFrames: collisionFrames,
				})
			}
		}
	}
}

// Important Events: OnTargetDied, OnTargetMoved (also emits on player move?), OnDendroCore
// Need to handle dendro core removal
func (c *char) bikeChargeAttackHook() {
	c.Core.Events.Subscribe(event.OnDendroCore, func(args ...interface{}) bool {
		// Ignore if not in bike state
		if c.armamentState != bike && !c.nightsoulState.HasBlessing() {
			return false
		}
		// If in bike state, recalculate gadget list
		g, ok := args[0].(combat.Gadget)
		if !ok {
			return false
		}
		if g.GadgetTyp() == combat.GadgetTypDendroCore {
			c.buildValidGadgetList()
		}

		return false
	}, "mavuika-bike-gadget-check")
}

func (c *char) bikeChargeAttackUnhook() {
	c.Core.Events.Unsubscribe(event.OnDendroCore, "mavuika-bike-gadget-check")
}

// Iterate through CA frames, starting at hitmark
func (c *char) CalculateValidCollisionFrames(durationCA int, collisionFrames [2]int) []int {
	validFrames := []int{}
	currentFrame := 35 // First hitmark occurs on frame 35
	// Check for and set 375f limit with prev CA time + new duration
	if durationCA+c.caState.CAtkFrames > bikeChargeAttackMaximumDuration {
		durationCA -= c.caState.CAtkFrames - bikeChargeAttackMaximumDuration
	}

	// Start at the hitmark
	timeSinceStart := c.Core.F - c.caState.StartFrame
	timeSinceLastHit := timeSinceStart - c.caState.LastHitF
	if timeSinceStart >= currentFrame {
		currentFrame = timeSinceStart
		if timeSinceLastHit < bikeChargeAttackICD {
			currentFrame += bikeChargeAttackICD - timeSinceLastHit
		}
	}
	totalFrames := currentFrame                              // Track the total frames elapsed
	currentFrame = currentFrame % bikeChargeAttackSpinFrames // Start current frame within spin cycle
	collisionStart := collisionFrames[0]
	collisionEnd := collisionFrames[1]

	c.Core.Log.NewEventBuildMsg(glog.LogSimEvent, c.Index,
		fmt.Sprintf("Duration: %d | Total frames: %d | Current frame: %d | Time since start: %d | CA frames: %d",
			durationCA, totalFrames, currentFrame, timeSinceStart, c.caState.CAtkFrames))

	for totalFrames <= (durationCA + c.caState.CAtkFrames) {
		checkValidFrame := -1
		// If the frame is outside the collision range, shift forward
		if collisionStart <= collisionEnd {
			if currentFrame > collisionEnd {
				currentFrame -= bikeChargeAttackSpinFrames
				totalFrames += collisionStart - currentFrame
				currentFrame = collisionStart
			} else if currentFrame < collisionStart {
				totalFrames += collisionStart - currentFrame
				currentFrame = collisionStart
			}
		} else {
			if currentFrame < collisionStart && currentFrame > collisionEnd {
				totalFrames += collisionStart - currentFrame
				currentFrame = collisionStart
			}
		}

		if collisionStart <= collisionEnd {
			if currentFrame >= collisionStart && currentFrame <= collisionEnd {
				checkValidFrame = totalFrames - timeSinceStart
			}
		} else {
			// Handle wrapping cases where collisionEnd is before collisionStart
			if currentFrame >= collisionStart || currentFrame <= collisionEnd {
				checkValidFrame = totalFrames - timeSinceStart
			}
		}
		if checkValidFrame >= 0 && checkValidFrame <= durationCA {
			validFrames = append(validFrames, checkValidFrame)
		}

		// Move forward by cooldownFrames, wrapping within spin animation length
		totalFrames += bikeChargeAttackICD
		currentFrame = (currentFrame + bikeChargeAttackICD) % bikeChargeAttackSpinFrames
	}

	return validFrames
}

// Calculate start and end frames for each spin during which target is within Mav hitbox
// Return false if target is not circle or has no overlap
func (c *char) BikeHitboxIntersectionAngles(v combat.Target, f []int) bool {
	enemyShape := v.Shape()
	var enemyRadius float64
	switch v := enemyShape.(type) {
	case *geometry.Circle:
		enemyRadius = v.Radius() // Rt
	default:
		return false
	}

	bikeInnerRadius := bikeChargeAttackSpinOffset - bikeChargeAttackHitboxRadius // Ri
	bikeOuterRadius := bikeChargeAttackSpinOffset + bikeChargeAttackHitboxRadius // Ro

	posDifference := v.Pos().Sub(c.Core.Combat.Player().Pos())
	enemyDistance := posDifference.Magnitude() // Dt

	// Check if no overlap
	if enemyDistance+enemyRadius < bikeInnerRadius || enemyDistance-enemyRadius > bikeOuterRadius {
		return false
	}

	// Target is always within hitbox range for the entire rotation
	if enemyRadius-enemyDistance > bikeInnerRadius {
		f[0] = 0
		f[1] = bikeChargeAttackSpinFrames
		return true
	} else if enemyDistance == 0 {
		return false
	}

	sumRadii := bikeChargeAttackHitboxRadius + enemyRadius
	cosThetaM := (bikeChargeAttackSpinOffset*bikeChargeAttackSpinOffset + enemyDistance*enemyDistance - sumRadii*sumRadii) /
		(2 * bikeChargeAttackSpinOffset * enemyDistance)

	if cosThetaM < -1 || cosThetaM > 1 {
		return false // No valid intersection??
	}

	enemyAngle := math.Atan2(posDifference.Y, posDifference.X) * (180 / math.Pi)
	thetaM := math.Acos(cosThetaM) * (180 / math.Pi)

	if enemyAngle < 0 {
		enemyAngle += 360
	}

	intersectAngleStart := enemyAngle + thetaM
	intersectAngleEnd := enemyAngle - thetaM

	// c.Core.Log.NewEventBuildMsg(glog.LogSimEvent, c.Index,
	// 	fmt.Sprintf("Intersection Angle Start: %.2f | End: %.2f | Enemy Radius: %.1f",
	// 		intersectAngleStart, intersectAngleEnd, enemyRadius))

	f[0] = c.ConvertAngleToFrame(intersectAngleStart, "start")
	f[1] = c.ConvertAngleToFrame(intersectAngleEnd, "end")

	return true
}

func (c *char) ConvertAngleToFrame(theta float64, s string) int {
	theta = math.Mod(theta+360, 360)

	var quadrant int
	var spinQuadrant int
	var accumulatedFrames int

	switch {
	case theta >= 270 || theta < 0:
		quadrant = 3
		spinQuadrant = 0
		accumulatedFrames = 0
	case theta >= 180:
		quadrant = 2
		spinQuadrant = 1
		accumulatedFrames = bikeSpinQuadrantFrames[0]
	case theta >= 90:
		quadrant = 1
		spinQuadrant = 2
		accumulatedFrames = bikeSpinQuadrantFrames[1] + bikeSpinQuadrantFrames[0]
	default:
		quadrant = 0
		spinQuadrant = 3
		accumulatedFrames = bikeSpinQuadrantFrames[2] + bikeSpinQuadrantFrames[1] + bikeSpinQuadrantFrames[0]
	}

	if accumulatedFrames > 0 {
		accumulatedFrames-- // Account for spin frame count starting at 0
	}

	// Calculate frame within quadrant
	quadrantStartAngle := float64(quadrant) * 90.0
	frameOffset := float64(bikeSpinQuadrantFrames[spinQuadrant]) + (theta-quadrantStartAngle)/bikeSpinQuadrantAngularVelocity[spinQuadrant]

	// c.Core.Log.NewEventBuildMsg(glog.LogSimEvent, c.Index,
	// 	fmt.Sprintf("%s: Theta: %.2f | quadrant: %d | AccumulatedFrames: %d | qStartAngle: %.2f | Offset: %.2f",
	// 		s, theta, quadrant, accumulatedFrames, quadrantStartAngle, frameOffset))

	return accumulatedFrames + int(math.Round(frameOffset))
}
