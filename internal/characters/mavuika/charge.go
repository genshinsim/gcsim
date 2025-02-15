package mavuika

import (
	"errors"
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var chargeFrames []int
var bikeChargeFrames []int
var bikeChargeFinalFrames []int
var bikeHittableEntityList []HittableEntity

// Minimum CA time before CAF anim is 50f
var bikeChargeAttackMinimumDuration = 50
var bikeChargeAttackStartupHitmark = 35

// Maximum CA time before CAF anim is 375f
var bikeChargeAttackMaximumDuration = 375
var bikeChargeFinalHitmark = 45

// TODO: Replicate frames 35-46 of the CA more accurately
// var bikeSpinInitialFrames = 11
// var bikeSpinInitialAngularVelocity = float64(-180 / 11)
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

	bikeChargeFrames = frames.InitAbilSlice(bikeChargeAttackMinimumDuration + bikeChargeFinalHitmark)
	bikeChargeFrames[action.ActionCharge] = 0
	bikeChargeFrames[action.ActionBurst] = 0
	bikeChargeFrames[action.ActionSkill] = 0
	bikeChargeFrames[action.ActionDash] = 0
	bikeChargeFrames[action.ActionJump] = 0
	bikeChargeFrames[action.ActionSwap] = 0

	bikeChargeFinalFrames = frames.InitAbilSlice(74) // CAF -> NA
	bikeChargeFinalFrames[action.ActionBurst] = bikeChargeFinalHitmark
	bikeChargeFinalFrames[action.ActionDash] = bikeChargeFinalHitmark
	bikeChargeFinalFrames[action.ActionJump] = bikeChargeFinalHitmark
	bikeChargeFinalFrames[action.ActionSwap] = bikeChargeFinalHitmark
	bikeChargeFinalFrames[action.ActionSkill] = bikeChargeFinalHitmark
}

// Charge state struct
type ChargeState struct {
	StartFrame      int
	cAtkFrames      int
	LastHit         map[int]int
	FacingDirection float64
	srcFrame        int
}

type HittableEntity struct {
	Entity     combat.Target
	isOneTick  bool   // Does entity get destroyed after a single maxHitCount?
	CollFrames [2]int // Frames of the CA spin on which collision happens
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if c.armamentState == bike && c.nightsoulState.HasBlessing() {
		return c.BikeCharge(p)
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
func (c *char) BikeCharge(p map[string]int) (action.Info, error) {
	// Parameters for tuning CA
	durationCA := p["hold"]
	final := p["final"]

	bikeHittableEntities := c.BuildBikeChargeAttackHittableTargetList()

	// Check if a continuing CA or new
	x := c.Core.Player.CurrentState()
	skippedWindupFrames := 0
	if x != action.ChargeAttackState || c.caState.StartFrame == 0 {
		c.caState = ChargeState{}
		c.caState.StartFrame = c.Core.F
		c.caState.LastHit = make(map[int]int)
		for _, t := range bikeHittableEntities {
			targetIndex := int(t.Entity.Key())
			c.caState.LastHit[targetIndex] = 0
		}
		c.Core.Tasks.Add(func() {
			c.bikeChargeAttackHook()
		}, bikeChargeAttackStartupHitmark-1)
		// TODO: Find better way of handling this
		// Currently the angle/hitbox tracking uses raw CA frames to determine position
		// Subtracting this at the wrong time can cause hits to get out of sync
		switch {
		case x == action.DashState:
			skippedWindupFrames = 15
			c.Core.Events.Emit(event.OnStateChange, action.NormalAttackState, action.NormalAttackState)
		case x == action.NormalAttackState || x == action.ChargeAttackState && c.caState.StartFrame == c.Core.F:
			skippedWindupFrames = 15
		case x == action.BurstState:
			burstFrames := c.Core.F - c.Core.Player.CurrentStateStart()
			if burstFrames >= 131 {
				c.Core.Events.Emit(event.OnStateChange, action.NormalAttackState, action.NormalAttackState)
			}
			skippedWindupFrames = 15
		}
	}

	nightSoulDuration := c.GetRemainingNightSoulDuration()
	isUseFinalHit := false
	isNightSoulExitQueued := false

	if final == 1 {
		return c.BikeChargeAttackFinal(0, skippedWindupFrames, false), nil
	}

	if durationCA > 0 && c.caState.cAtkFrames > 0 {
		// Limit hold duration to 1 full spin
		if durationCA > bikeChargeAttackSpinFrames {
			// Default max hold time is 6.25s/375f
			if (c.caState.cAtkFrames + durationCA) <= bikeChargeAttackMaximumDuration {
				durationCA = bikeChargeAttackSpinFrames
			} else {
				durationCA = bikeChargeAttackMaximumDuration - c.caState.cAtkFrames
			}
		} else if durationCA > nightSoulDuration {
			durationCA = nightSoulDuration
		}
		// Hold CA logic
		c.HoldBikeChargeAttack(durationCA, skippedWindupFrames, bikeHittableEntities)
	} else {
		// Default to single CA if nothing specified
		hasNoValidTarget, ai, err := c.HasValidTargetCheck(bikeHittableEntities)
		if hasNoValidTarget {
			return ai, err
		}
		durationCA = c.CountBikeChargeAttack(1, skippedWindupFrames, bikeHittableEntities, nightSoulDuration)
	}

	// Add any existing CA frames
	c.caState.cAtkFrames += durationCA
	durationCA -= skippedWindupFrames

	if durationCA >= nightSoulDuration {
		isUseFinalHit = true
		isNightSoulExitQueued = true
		c.QueueCharTask(func() {
			c.exitNightsoul()
		}, nightSoulDuration)
	}

	if isUseFinalHit {
		return c.BikeChargeAttackFinal(durationCA, skippedWindupFrames, isNightSoulExitQueued), nil
	}

	c.caState.srcFrame = c.Core.F
	src := c.caState.srcFrame

	// Start queue CAF for invalid actions
	// Check if bike angle is in spot where CAF has delay, 20f window (used for CAF queue)
	currentBikeSpinFrame := c.caState.cAtkFrames % bikeChargeAttackSpinFrames
	newMinSpinDuration := GetCAFDelay(currentBikeSpinFrame)

	c.QueueCharTask(func() {
		if c.caState.srcFrame != src {
			return
		}
		c.HoldBikeChargeAttack(newMinSpinDuration, skippedWindupFrames, bikeHittableEntities)
	}, durationCA+1)

	c.QueueCharTask(func() {
		if c.caState.srcFrame != src {
			return
		}
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
		c.caState = ChargeState{}
		c.bikeChargeAttackUnhook()
	}, durationCA+newMinSpinDuration+bikeChargeFinalHitmark)

	return action.Info{
		Frames: func(next action.Action) int {
			f := bikeChargeFrames[next]

			if f == 0 {
				f = durationCA
			} else {
				f = durationCA + newMinSpinDuration + bikeChargeFinalFrames[next]
			}
			return f
		},
		AnimationLength: durationCA + newMinSpinDuration + bikeChargeFinalFrames[action.InvalidAction],
		CanQueueAfter:   durationCA,
		State:           action.ChargeAttackState,
		OnRemoved: func(next action.AnimationState) {
			if next != action.ChargeAttackState {
				c.caState = ChargeState{}
				c.bikeChargeAttackUnhook()
			}
		},
	}, nil
}

// For given CA length, calculate hits on each target in hittable list
func (c *char) HoldBikeChargeAttack(cAtkFrames int, skippedWindupFrames int, hittableEntities []HittableEntity) {
	for i := 0; i < len(hittableEntities); i++ {
		t := hittableEntities[i]
		enemyID := int(t.Entity.Key())
		lastHitFrame := c.caState.LastHit[enemyID]
		newLastHitFrame := -1

		if t.isOneTick && lastHitFrame > 0 {
			continue
		}

		// First 11f of CA are a bit inaccurate, should maxHitCount further left
		hitFrames := c.CalculateValidCollisionFrames(cAtkFrames, t.CollFrames, lastHitFrame)

		if len(hitFrames) > 0 {
			for _, f := range hitFrames {
				c.Core.Tasks.Add(func() {
					ai := c.GetBikeChargeAttackAttackInfo()
					c.Core.QueueAttack(ai, combat.NewSingleTargetHit(t.Entity.Key()), 0, 0)
				}, f-skippedWindupFrames)
				newLastHitFrame = f
			}
		}
		if newLastHitFrame >= 0 {
			c.caState.LastHit[enemyID] += newLastHitFrame + (c.caState.cAtkFrames - lastHitFrame)
		}
	}
}

// For given maxHitCount count, calculate maxHitCount timings on targets and return CA duration
func (c *char) CountBikeChargeAttack(maxHitCount int, skippedWindupFrames int, hittableEntities []HittableEntity, nsDur int) int {
	// Return remaining CA time between nightsoul duration and max CA duration for attempting hit
	dur := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}(nsDur, bikeChargeAttackMaximumDuration-c.caState.cAtkFrames)
	hitCounter := 0

	for i := 0; i < len(hittableEntities); i++ {
		t := hittableEntities[i]
		if t.Entity != c.Core.Combat.PrimaryTarget() {
			continue
		}

		enemyID := int(t.Entity.Key())
		lastHitFrame := c.caState.LastHit[enemyID]

		if t.isOneTick && lastHitFrame > 0 {
			continue
		}

		// First 11f of CA are a bit inaccurate, should start further left and sweep faster
		hitFrames := c.CalculateValidCollisionFrames(dur, t.CollFrames, lastHitFrame)

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

	for i := 0; i < len(hittableEntities); i++ {
		t := hittableEntities[i]
		enemyID := int(t.Entity.Key())
		lastHitFrame := c.caState.LastHit[enemyID]
		newLastHitFrame := -1

		// First 11f of CA are a bit inaccurate, should maxHitCount further left
		hitFrames := c.CalculateValidCollisionFrames(dur, t.CollFrames, lastHitFrame)

		if len(hitFrames) > 0 {
			for _, f := range hitFrames {
				c.Core.Tasks.Add(func() {
					ai := c.GetBikeChargeAttackAttackInfo()
					c.Core.QueueAttack(ai, combat.NewSingleTargetHit(t.Entity.Key()), 0, 0)
				}, f-skippedWindupFrames)
				newLastHitFrame = f
			}
		}
		// Used when the CA started between hits (Usually for secondary+ targets)
		if newLastHitFrame >= 0 {
			c.caState.LastHit[enemyID] += newLastHitFrame + (c.caState.cAtkFrames - lastHitFrame)
		}
	}
	return dur
}

func (c *char) BikeChargeAttackFinal(caFrames, skippedWindupFrames int, hasQueuedExitNightSoul bool) action.Info {
	var adjustedBikeChargeFinalHitmark = bikeChargeFinalHitmark + caFrames
	bikeChargeAttackElapsedTime := c.caState.cAtkFrames + caFrames
	// Check if bike angle is in spot where CAF has delay, 20f window
	currentBikeSpinFrame := bikeChargeAttackElapsedTime % bikeChargeAttackSpinFrames
	newMinSpinDuration := GetCAFDelay(currentBikeSpinFrame)

	caFrames += newMinSpinDuration
	adjustedBikeChargeFinalHitmark += caFrames
	bikeHittableEntities := c.BuildBikeChargeAttackHittableTargetList()
	c.HoldBikeChargeAttack(newMinSpinDuration, skippedWindupFrames, bikeHittableEntities)

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

	nightSoulDuration := c.GetRemainingNightSoulDuration()
	if nightSoulDuration <= adjustedBikeChargeFinalHitmark {
		// Exiting at hitmark to account for dash cancel
		c.QueueCharTask(func() {
			c.exitBike()
		}, adjustedBikeChargeFinalHitmark)

		if !hasQueuedExitNightSoul {
			c.QueueCharTask(func() {
				c.exitNightsoul()
			}, nightSoulDuration)
		}
	}

	c.Core.Tasks.Add(func() {
		c.caState = ChargeState{}
		c.bikeChargeAttackUnhook()
	}, caFrames)

	return action.Info{
		Frames:          func(next action.Action) int { return bikeChargeFinalFrames[next] + caFrames },
		AnimationLength: bikeChargeFinalFrames[action.InvalidAction] + caFrames,
		CanQueueAfter:   bikeChargeFinalFrames[action.ActionDash] + caFrames,
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

// CA NightSoul consumption is 11/s, with skill.go function reducing this every 6f
func (c *char) GetRemainingNightSoulDuration() int {
	curPoints := c.nightsoulState.Points()
	// This will give up to +5f of duration, not synced with skill.go
	nsDur := int(math.Ceil(curPoints / 1.1))
	nsDur *= 6
	if c.StatusIsActive(burstKey) {
		nsDur += c.StatusDuration(burstKey)
	}

	return nsDur
}

// No idea why this is asymmetrical across the 0 degree mark, but testing showed a big jump across that line
func GetCAFDelay(currentBikeSpinFrame int) int {
	newMinSpinDuration := 0

	if currentBikeSpinFrame < 10 {
		newMinSpinDuration = 10 - currentBikeSpinFrame
	} else if currentBikeSpinFrame >= 35 {
		newMinSpinDuration = 50 - currentBikeSpinFrame
	}
	return newMinSpinDuration
}

func (c *char) BuildBikeChargeAttackHittableTargetList() []HittableEntity {
	return append(c.buildValidTargetList(), c.buildValidGadgetList()...)
}

func (c *char) buildValidTargetList() []HittableEntity {
	enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 8), nil)
	hittableEnemies := []HittableEntity{}
	for _, v := range enemies {
		if v == nil {
			continue
		}
		// Calculate start and ending frames for collision
		collisionFrames := [2]int{-1, -1}
		facingDirection := 0.0
		if c.caState.cAtkFrames == 0 {
			facingDirection = c.DirectionOffsetToPrimaryTarget()
			c.caState.FacingDirection = facingDirection
		} else {
			facingDirection = c.caState.FacingDirection
		}
		isIntersecting := c.BikeHitboxIntersectionAngles(v, collisionFrames[:], facingDirection)

		if isIntersecting {
			hittableEnemies = append(hittableEnemies, HittableEntity{
				Entity:     combat.Target(v),
				isOneTick:  false,
				CollFrames: collisionFrames,
			})
		}
	}
	// c.Core.Log.NewEventBuildMsg(glog.LogSimEvent, c.Index,
	// 	fmt.Sprintf("There are %v valid targets", len(bikeHittableEntities)))
	return hittableEnemies
}

// Gadgets are gonna be problematic
func (c *char) buildValidGadgetList() []HittableEntity {
	gadgets := c.Core.Combat.GadgetsWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 8), nil)
	hittableGadgets := []HittableEntity{}
	for _, g := range gadgets {
		if g == nil {
			continue
		}
		if g.GadgetTyp() == combat.GadgetTypDendroCore {
			// Calculate start and ending frames for collision
			hittableGadget, isHittable := c.AddNewGadgetToHittableList(g)
			if isHittable {
				hittableGadgets = append(hittableGadgets, hittableGadget)
			}
		}
	}
	return hittableGadgets
}

func (c *char) AddNewGadgetToHittableList(v combat.Gadget) (HittableEntity, bool) {
	collisionFrames := [2]int{-1, -1}
	facingDirection := 0.0
	if c.caState.cAtkFrames == 0 {
		facingDirection = c.DirectionOffsetToPrimaryTarget()
		c.caState.FacingDirection = facingDirection
	} else {
		facingDirection = c.caState.FacingDirection
	}
	isIntersecting := c.BikeHitboxIntersectionAngles(v, collisionFrames[:], facingDirection)
	newGadget := HittableEntity{}

	if isIntersecting {
		newGadget = HittableEntity{
			Entity:     combat.Target(v),
			isOneTick:  true,
			CollFrames: collisionFrames,
		}
	}
	return newGadget, isIntersecting
}

func (c *char) HasValidTargetCheck(bikeHittableEntities []HittableEntity) (bool, action.Info, error) {
	isTargetForCountsHittable := false
	if len(bikeHittableEntities) == 0 {
		c.SetHittableEntityList(bikeHittableEntities)
		return true, action.Info{}, errors.New("no valid targets within flamestrider area")
	}
	for _, t := range bikeHittableEntities {
		if t.Entity == c.Core.Combat.PrimaryTarget() {
			isTargetForCountsHittable = true
			break
		}
	}
	if !isTargetForCountsHittable {
		return true, action.Info{}, errors.New("primary target is not within flamestrider area")
	}
	return false, action.Info{}, nil
}

// Important Events: OnTargetDied, OnTargetMoved (also emits on player move?), OnDendroCore
func (c *char) bikeChargeAttackHook() {
	c.Core.Events.Subscribe(event.OnDendroCore, func(args ...interface{}) bool {
		// Ignore if not in bike state
		if c.armamentState != bike && !c.nightsoulState.HasBlessing() {
			return false
		}
		// If in bike state, add gadget to target list if it can be hit
		g, ok := args[0].(combat.Gadget)
		if !ok {
			return false
		}
		if g.GadgetTyp() == combat.GadgetTypDendroCore {
			hittableGadgets := c.GetHittableEntityList()
			// Might not be necessary to add to list?
			hittableGadget, isHittable := c.AddNewGadgetToHittableList(g)
			if isHittable {
				hittableGadgets = append(hittableGadgets, hittableGadget)
				remainingCADuration := c.caState.cAtkFrames - (c.Core.F - c.caState.StartFrame)
				x := c.Core.Player.CurrentState()
				skippedWindupFrames := 0
				if x == action.DashState || x == action.NormalAttackState {
					skippedWindupFrames = 15
				}
				hitFrames := c.CalculateValidCollisionFrames(remainingCADuration, hittableGadget.CollFrames, 0)
				if len(hitFrames) > 0 {
					for _, f := range hitFrames {
						c.Core.Tasks.Add(func() {
							ai := c.GetBikeChargeAttackAttackInfo()
							c.Core.QueueAttack(ai, combat.NewSingleTargetHit(hittableGadget.Entity.Key()), 0, 0)
						}, f-skippedWindupFrames)
					}
				}
				// Frame doesn't really matter as long as > 0
				c.caState.LastHit[int(g.Key())] += c.Core.F
			}
		}

		return false
	}, "mavuika-bike-gadget-check")
}

func (c *char) bikeChargeAttackUnhook() {
	c.Core.Events.Unsubscribe(event.OnDendroCore, "mavuika-bike-gadget-check")
}

func (*char) SetHittableEntityList(bikeHittableEntities []HittableEntity) {
	bikeHittableEntityList = bikeHittableEntities
}

func (*char) GetHittableEntityList() []HittableEntity {
	return bikeHittableEntityList
}

// Iterate through CA frames, starting at hitmark
func (c *char) CalculateValidCollisionFrames(durationCA int, collisionFrames [2]int, lastHitFrame int) []int {
	validFrames := []int{}
	currentFrame := 35 // First hitmark occurs on frame 35

	// Start at the hitmark
	timeSinceStart := c.caState.cAtkFrames
	timeSinceLastHit := timeSinceStart - lastHitFrame
	if timeSinceStart >= currentFrame {
		currentFrame = timeSinceStart
		if timeSinceLastHit < bikeChargeAttackICD {
			currentFrame += bikeChargeAttackICD - timeSinceLastHit

		}
	}
	totalFrames := currentFrame                // Track the total frames elapsed
	currentFrame %= bikeChargeAttackSpinFrames // Start current frame within spin cycle

	collisionStart := collisionFrames[0]
	collisionEnd := collisionFrames[1]

	for totalFrames <= (durationCA + c.caState.cAtkFrames) {
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
func (c *char) BikeHitboxIntersectionAngles(v combat.Target, f []int, offsetAngle float64) bool {
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
	if enemyDistance+enemyRadius <= bikeInnerRadius || enemyDistance-enemyRadius >= bikeOuterRadius {
		return false
	}

	// Target is always within hitbox range for the entire rotation
	if enemyRadius-enemyDistance > bikeInnerRadius {
		f[0] = 0
		f[1] = bikeChargeAttackSpinFrames
		return true
	}

	sumRadii := bikeChargeAttackHitboxRadius + enemyRadius
	cosThetaM := (bikeChargeAttackSpinOffset*bikeChargeAttackSpinOffset + enemyDistance*enemyDistance - sumRadii*sumRadii) /
		(2 * bikeChargeAttackSpinOffset * enemyDistance)

	enemyAngle := math.Atan2(posDifference.Y, posDifference.X) * (180 / math.Pi)
	thetaM := math.Acos(cosThetaM) * (180 / math.Pi)

	enemyAngle = math.Mod(enemyAngle-offsetAngle+360, 360)

	intersectAngleStart := enemyAngle + thetaM
	intersectAngleEnd := enemyAngle - thetaM

	f[0] = c.ConvertAngleToFrame(intersectAngleStart)
	f[1] = c.ConvertAngleToFrame(intersectAngleEnd)

	return true
}

func (c *char) DirectionOffsetToPrimaryTarget() float64 {
	var enemyDirection = geometry.CalcDirection(c.Core.Combat.Player().Pos(), c.Core.Combat.PrimaryTarget().Pos())
	if enemyDirection == geometry.DefaultDirection() {
		return 0
	}

	angleToTarget := math.Atan2(enemyDirection.X, enemyDirection.Y) * (180 / math.Pi)
	offsetAngle := 360 - angleToTarget

	return math.Mod(offsetAngle+360, 360)
}

func (c *char) ConvertAngleToFrame(theta float64) int {
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

	return accumulatedFrames + int(math.Round(frameOffset))
}
