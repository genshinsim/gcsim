package alhaitham

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

var lowPlungeFramesAL []int

const lowPlungeHitmarkAL = 38

const lowPlungeHitmarkXY = 46
const highPlungeHitmarkXY = 48
const collisionHitmarkXY = lowPlungeHitmarkXY - 6

const lowPlungePoiseDMG = 100.0
const lowPlungeRadius = 3.0

const highPlungePoiseDMG = 150.0
const highPlungeRadius = 5.0

var highPlungeFramesXY []int
var lowPlungeFramesXY []int

func init() {
	lowPlungeFramesAL = frames.InitAbilSlice(70)
	lowPlungeFramesAL[action.ActionAttack] = 49
	lowPlungeFramesAL[action.ActionSkill] = 50
	lowPlungeFramesAL[action.ActionBurst] = 50
	lowPlungeFramesAL[action.ActionDash] = 40
	lowPlungeFramesAL[action.ActionSwap] = 58

	// low_plunge -> x
	lowPlungeFramesXY = frames.InitAbilSlice(75)
	lowPlungeFramesXY[action.ActionAttack] = 53
	lowPlungeFramesXY[action.ActionSkill] = 54
	lowPlungeFramesXY[action.ActionBurst] = 55
	lowPlungeFramesXY[action.ActionDash] = lowPlungeHitmarkXY
	lowPlungeFramesXY[action.ActionJump] = 73
	lowPlungeFramesXY[action.ActionSwap] = 61

	// high_plunge -> x
	highPlungeFramesXY = frames.InitAbilSlice(77)
	highPlungeFramesXY[action.ActionAttack] = 56
	highPlungeFramesXY[action.ActionSkill] = 56
	highPlungeFramesXY[action.ActionBurst] = 56
	highPlungeFramesXY[action.ActionDash] = highPlungeHitmarkXY
	highPlungeFramesXY[action.ActionJump] = 76
	highPlungeFramesXY[action.ActionSwap] = 64
}

func (c *char) LowPlungeAttack(p map[string]int) (action.Info, error) {
	defer c.Core.Player.SetAirborne(player.Grounded)
	// last action hold skill
	if c.Core.Player.LastAction.Type == action.ActionSkill &&
		c.Core.Player.LastAction.Param["hold"] == 1 {
		return c.lowPlungeAl(p), nil
	}

	switch c.Core.Player.Airborne() {
	case player.AirborneXianyun:
		return c.lowPlungeXY(p), nil
	default:
		return action.Info{}, errors.New("low_plunge can only be used while airborne or after hold skill")
	}
}

func (c *char) lowPlungeAl(p map[string]int) action.Info {
	short := p["short"]
	skip := 0
	if short > 0 {
		skip = 20
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Low Plunge Attack",
		AttackTag:  attacks.AttackTagPlunge,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   lowPlungePoiseDMG,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       lowPlunge[c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, 3),
		lowPlungeHitmarkAL-skip,
		lowPlungeHitmarkAL-skip,
		c.makeA1CB(), // A1 adds a stack before the mirror count for the Projection Attack is determined
		c.projectionAttack,
	)

	return action.Info{
		Frames:          func(next action.Action) int { return lowPlungeFramesAL[next] - skip },
		AnimationLength: lowPlungeFramesAL[action.InvalidAction] - skip,
		CanQueueAfter:   lowPlungeFramesAL[action.ActionDash] - skip,
		State:           action.PlungeAttackState,
	}
}

func (c *char) lowPlungeXY(p map[string]int) action.Info {
	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not collision hit
	}

	if collision > 0 {
		c.plungeCollision(collisionHitmarkXY)
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
		lowPlungeHitmarkXY,
		lowPlungeHitmarkXY,
		c.makeA1CB(), // A1 adds a stack before the mirror count for the Projection Attack is determined
		c.projectionAttack,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(lowPlungeFramesAL),
		AnimationLength: lowPlungeFramesAL[action.InvalidAction],
		CanQueueAfter:   lowPlungeFramesAL[action.ActionDash],
		State:           action.PlungeAttackState,
	}
}

// High Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *char) HighPlungeAttack(p map[string]int) (action.Info, error) {
	defer c.Core.Player.SetAirborne(player.Grounded)
	switch c.Core.Player.Airborne() {
	case player.AirborneXianyun:
		return c.highPlungeXY(p), nil
	default:
		return action.Info{}, errors.New("high_plunge can only be used while airborne")
	}
}

func (c *char) highPlungeXY(p map[string]int) action.Info {
	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not collision hit
	}

	if collision > 0 {
		c.plungeCollision(collisionHitmarkXY)
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
		highPlungeHitmarkXY,
		highPlungeHitmarkXY,
		c.makeA1CB(), // A1 adds a stack before the mirror count for the Projection Attack is determined
		c.projectionAttack,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(highPlungeFramesXY),
		AnimationLength: highPlungeFramesXY[action.InvalidAction],
		CanQueueAfter:   highPlungeFramesXY[action.ActionDash],
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
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, 1),
		delay,
		delay,
		c.makeA1CB(), // A1 adds a stack before the mirror count for the Projection Attack is determined
		c.projectionAttack,
	)
}
