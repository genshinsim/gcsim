package raiden

import (
	"errors"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var highPlungeFrames []int
var lowPlungeFrames []int

var highPlungeFramesB []int
var lowPlungeFramesB []int

const lowPlungeHitmark = 44 + 3
const highPlungeHitmark = 46 + 3
const collisionHitmark = lowPlungeHitmark - 6

const lowPlungeHitmarkB = 46 + 3
const highPlungeHitmarkB = 48 + 3
const collisionHitmarkB = lowPlungeHitmarkB - 6

const lowPlungePoiseDMG = 100.0
const lowPlungeRadius = 3.0

const highPlungePoiseDMG = 150.0
const highPlungeRadius = 5.0

func init() {
	// low_plunge -> x
	lowPlungeFrames = frames.InitAbilSlice(96)
	lowPlungeFrames[action.ActionAttack] = 60
	lowPlungeFrames[action.ActionSkill] = 59
	lowPlungeFrames[action.ActionBurst] = 60
	lowPlungeFrames[action.ActionDash] = 80
	lowPlungeFrames[action.ActionJump] = 80
	lowPlungeFrames[action.ActionSwap] = 62

	// high_plunge -> x
	highPlungeFrames = frames.InitAbilSlice(99)
	highPlungeFrames[action.ActionAttack] = 61
	highPlungeFrames[action.ActionSkill] = 60
	highPlungeFrames[action.ActionBurst] = 62
	highPlungeFrames[action.ActionDash] = 84
	highPlungeFrames[action.ActionJump] = 84
	highPlungeFrames[action.ActionSwap] = 65

	// low_plunge -> x
	lowPlungeFramesB = frames.InitAbilSlice(85)
	lowPlungeFramesB[action.ActionAttack] = 58
	lowPlungeFramesB[action.ActionSkill] = 59
	lowPlungeFramesB[action.ActionBurst] = 59 // Assuming same as skill
	lowPlungeFramesB[action.ActionDash] = 78
	lowPlungeFramesB[action.ActionJump] = 81
	lowPlungeFramesB[action.ActionSwap] = 60

	// high_plunge -> x
	highPlungeFramesB = frames.InitAbilSlice(88)
	highPlungeFramesB[action.ActionAttack] = 60
	highPlungeFramesB[action.ActionSkill] = 60
	highPlungeFramesB[action.ActionBurst] = 60 // Assuming same as skill
	highPlungeFramesB[action.ActionDash] = 83
	highPlungeFramesB[action.ActionJump] = 81
	highPlungeFramesB[action.ActionSwap] = 61
}

// Low Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *char) LowPlungeAttack(p map[string]int) (action.Info, error) {
	defer c.Core.Player.SetAirborne(player.Grounded)
	switch c.Core.Player.Airborne() {
	case player.AirborneXianyun:
		if c.StatusIsActive(BurstKey) {
			return c.lowPlungeBXY(p)
		}
		return c.lowPlungeXY(p)
	default:
		return action.Info{}, errors.New("low_plunge can only be used while airborne")
	}
}

func (c *char) lowPlungeXY(p map[string]int) (action.Info, error) {
	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not collision hit
	}

	if collision > 0 {
		c.plungeCollision(collisionHitmark)
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Low Plunge",
		AttackTag:  attacks.AttackTagPlunge,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   lowPlungePoiseDMG,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       lowPlunge[c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, lowPlungeRadius),
		lowPlungeHitmark,
		lowPlungeHitmark,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(lowPlungeFrames),
		AnimationLength: lowPlungeFrames[action.InvalidAction],
		CanQueueAfter:   lowPlungeFrames[action.ActionSkill],
		State:           action.PlungeAttackState,
	}, nil
}

func (c *char) lowPlungeBXY(p map[string]int) (action.Info, error) {
	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not collision hit
	}

	if collision > 0 {
		c.plungeCollisionB(collisionHitmarkB)
	}

	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Low Plunge (Q)",
		AttackTag:      attacks.AttackTagElementalBurst,
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeBlunt,
		PoiseDMG:       lowPlungePoiseDMG,
		Element:        attributes.Electro,
		IgnoreInfusion: true,
		Durability:     25,
		Mult:           lowPlunge[c.TalentLvlBurst()] + resolveBonus[c.TalentLvlBurst()]*c.stacksConsumed,
	}
	if c.Base.Cons >= 2 {
		ai.IgnoreDefPercent = .6
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, lowPlungeRadius),
		lowPlungeHitmarkB,
		lowPlungeHitmarkB,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(lowPlungeFramesB),
		AnimationLength: lowPlungeFramesB[action.InvalidAction],
		CanQueueAfter:   lowPlungeFramesB[action.ActionAttack],
		State:           action.PlungeAttackState,
	}, nil
}

// High Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *char) HighPlungeAttack(p map[string]int) (action.Info, error) {
	defer c.Core.Player.SetAirborne(player.Grounded)
	switch c.Core.Player.Airborne() {
	case player.AirborneXianyun:
		if c.StatusIsActive(BurstKey) {
			return c.highPlungeBXY(p)
		}
		return c.highPlungeXY(p)
	default:
		return action.Info{}, errors.New("high_plunge can only be used while airborne")
	}
}

func (c *char) highPlungeXY(p map[string]int) (action.Info, error) {
	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not collision hit
	}

	if collision > 0 {
		c.plungeCollision(collisionHitmark)
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "High Plunge",
		AttackTag:  attacks.AttackTagPlunge,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   highPlungePoiseDMG,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       highPlunge[c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, highPlungeRadius),
		highPlungeHitmark,
		highPlungeHitmark,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(highPlungeFrames),
		AnimationLength: highPlungeFrames[action.InvalidAction],
		CanQueueAfter:   highPlungeFrames[action.ActionAttack],
		State:           action.PlungeAttackState,
	}, nil
}

func (c *char) highPlungeBXY(p map[string]int) (action.Info, error) {
	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not collision hit
	}

	if collision > 0 {
		c.plungeCollisionB(collisionHitmarkB)
	}

	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "High Plunge (Q)",
		AttackTag:      attacks.AttackTagElementalBurst,
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeBlunt,
		PoiseDMG:       highPlungePoiseDMG,
		Element:        attributes.Electro,
		IgnoreInfusion: true,
		Durability:     25,
		Mult:           highPlunge[c.TalentLvlBurst()] + resolveBonus[c.TalentLvlBurst()]*c.stacksConsumed,
	}
	if c.Base.Cons >= 2 {
		ai.IgnoreDefPercent = .6
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, highPlungeRadius),
		highPlungeHitmarkB,
		highPlungeHitmarkB,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(highPlungeFramesB),
		AnimationLength: highPlungeFramesB[action.InvalidAction],
		CanQueueAfter:   highPlungeFramesB[action.ActionSkill],
		State:           action.PlungeAttackState,
	}, nil
}

// Plunge normal falling attack damage queue generator
// Standard - Always part of high/low plunge attacks
func (c *char) plungeCollision(delay int) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Plunge Collision",
		AttackTag:  attacks.AttackTagPlunge,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Physical,
		Durability: 0,
		Mult:       collision[c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, 1), delay, delay)
}

// Plunge normal falling attack damage queue generator
// Standard - Always part of high/low plunge attacks
func (c *char) plungeCollisionB(delay int) {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Plunge Collision",
		AttackTag:      attacks.AttackTagElementalBurst,
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeSlash,
		Element:        attributes.Electro,
		IgnoreInfusion: true,
		Durability:     0,
		Mult:           collision[c.TalentLvlBurst()] + resolveBonus[c.TalentLvlBurst()]*c.stacksConsumed,
	}
	if c.Base.Cons >= 2 {
		ai.IgnoreDefPercent = .6
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, 1), delay, delay)
}
