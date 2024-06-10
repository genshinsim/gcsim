package kazuha

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

var plungePressFrames []int
var plungeHoldFrames []int

// a1 is 1 frame before this
// collision is 6 frame before this
const plungePressHitmark = 36
const plungeHoldHitmark = 41

var highPlungeFrames []int
var lowPlungeFrames []int

const lowPlungeHitmark = 46
const highPlungeHitmark = 47
const collisionHitmark = lowPlungeHitmark - 6

const lowPlungePoiseDMG = 100.0
const lowPlungeRadius = 3.0

const highPlungePoiseDMG = 150.0
const highPlungeRadius = 5.0

// TODO: missing plunge -> skill
func init() {
	// skill (press) -> high plunge -> x
	plungePressFrames = frames.InitAbilSlice(55) // max
	plungePressFrames[action.ActionDash] = 43
	plungePressFrames[action.ActionJump] = 50
	plungePressFrames[action.ActionSwap] = 50

	// skill (hold) -> high plunge -> x
	plungeHoldFrames = frames.InitAbilSlice(61) // max
	plungeHoldFrames[action.ActionSkill] = 60   // uses burst frames
	plungeHoldFrames[action.ActionBurst] = 60
	plungeHoldFrames[action.ActionDash] = 48
	plungeHoldFrames[action.ActionJump] = 55
	plungeHoldFrames[action.ActionSwap] = 54

	// low_plunge -> x
	lowPlungeFrames = frames.InitAbilSlice(73)
	lowPlungeFrames[action.ActionAttack] = 52
	lowPlungeFrames[action.ActionSkill] = 52
	lowPlungeFrames[action.ActionBurst] = 51
	lowPlungeFrames[action.ActionDash] = lowPlungeHitmark
	lowPlungeFrames[action.ActionJump] = 69
	lowPlungeFrames[action.ActionSwap] = 53

	// high_plunge -> x
	highPlungeFrames = frames.InitAbilSlice(73)
	highPlungeFrames[action.ActionAttack] = 54
	highPlungeFrames[action.ActionSkill] = 53
	highPlungeFrames[action.ActionBurst] = 53
	highPlungeFrames[action.ActionDash] = highPlungeHitmark
	highPlungeFrames[action.ActionJump] = 69
	highPlungeFrames[action.ActionSwap] = 55
}

// Low Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *char) LowPlungeAttack(p map[string]int) (action.Info, error) {
	defer c.Core.Player.SetAirborne(player.Grounded)
	if c.Core.Player.LastAction.Type == action.ActionSkill {
		return action.Info{}, errors.New("cannot low_plunge after skill")
	}

	switch c.Core.Player.Airborne() {
	case player.AirborneXianyun:
		return c.lowPlungeXY(p), nil
	default:
		return action.Info{}, errors.New("low_plunge can only be used while airborne")
	}
}

func (c *char) lowPlungeXY(p map[string]int) action.Info {
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
		CanQueueAfter:   lowPlungeFrames[action.ActionDash],
		State:           action.PlungeAttackState,
	}
}

func (c *char) HighPlungeAttack(p map[string]int) (action.Info, error) {
	defer c.Core.Player.SetAirborne(player.Grounded)
	if c.Core.Player.LastAction.Type == action.ActionSkill {
		return c.skillPlunge(p)
	}

	switch c.Core.Player.Airborne() {
	case player.AirborneXianyun:
		return c.highPlungeXY(p), nil
	default:
		return action.Info{}, errors.New("high_plunge can only be used while airborne")
	}
}

func (c *char) skillPlunge(p map[string]int) (action.Info, error) {
	// last action must be skill without glide cancel
	if c.Core.Player.LastAction.Param["glide_cancel"] != 0 {
		return action.Info{}, errors.New("only plunge after skill without glide cancel")
	}

	act := action.Info{
		State: action.PlungeAttackState,
	}

	//TODO: is this accurate?? these should be the hitmarks
	var hitmark int
	if c.Core.Player.LastAction.Param["hold"] == 0 {
		hitmark = plungePressHitmark
		act.Frames = frames.NewAbilFunc(plungePressFrames)
		act.AnimationLength = plungePressFrames[action.InvalidAction]
		act.CanQueueAfter = plungePressFrames[action.ActionDash] // earliest cancel
	} else {
		hitmark = plungeHoldHitmark
		act.Frames = frames.NewAbilFunc(plungeHoldFrames)
		act.AnimationLength = plungeHoldFrames[action.InvalidAction]
		act.CanQueueAfter = plungeHoldFrames[action.ActionDash] // earliest cancel
	}
	collisionParam, ok := p["collision"]
	if !ok {
		collisionParam = 0 // Whether or not collision hit
	}

	if collisionParam > 0 {
		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Plunge Collision",
			AttackTag:      attacks.AttackTagPlunge,
			ICDTag:         attacks.ICDTagNone,
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypeSlash,
			Element:        attributes.Anemo,
			IgnoreInfusion: true,
			Durability:     0,
			Mult:           collision[c.TalentLvlAttack()],
		}
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, 1), hitmark-6, hitmark-6)
	}

	// aoe dmg
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "High Plunge",
		AttackTag:      attacks.AttackTagPlunge,
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeBlunt,
		PoiseDMG:       150,
		Element:        attributes.Anemo,
		Durability:     25,
		Mult:           highPlunge[c.TalentLvlAttack()],
		IgnoreInfusion: true,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 0.5}, 4.5),
		hitmark,
		hitmark,
	)

	// a1 if applies
	if c.a1Absorb != attributes.NoElement {
		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Kazuha A1",
			AttackTag:      attacks.AttackTagPlunge,
			ICDTag:         attacks.ICDTagNone,
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypeBlunt,
			PoiseDMG:       20,
			Element:        c.a1Absorb,
			Durability:     25,
			Mult:           2,
			IgnoreInfusion: true,
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 0.5}, 4.5),
			hitmark-1,
			hitmark-1,
		)
		c.a1Absorb = attributes.NoElement
	}

	return act, nil
}

func (c *char) highPlungeXY(p map[string]int) action.Info {
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
		CanQueueAfter:   highPlungeFrames[action.ActionDash],
		State:           action.PlungeAttackState,
	}
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
