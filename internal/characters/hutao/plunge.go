package hutao

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

const lowPlungeHitmark = 41 + 3
const highPlungeHitmark = 42 + 3
const collisionHitmark = lowPlungeHitmark - 6

const lowPlungeHitmarkB = 42 + 3
const highPlungeHitmarkB = 43 + 3
const collisionHitmarkB = lowPlungeHitmarkB - 6

const lowPlungePoiseDMG = 100.0
const lowPlungeRadius = 3.0

const highPlungePoiseDMG = 150.0
const highPlungeRadius = 5.0

func init() {
	// low_plunge -> x
	lowPlungeFrames = frames.InitAbilSlice(75)
	lowPlungeFrames[action.ActionAttack] = 58
	lowPlungeFrames[action.ActionSkill] = 59
	lowPlungeFrames[action.ActionBurst] = 59
	lowPlungeFrames[action.ActionJump] = 74
	lowPlungeFrames[action.ActionSwap] = 61

	// high_plunge -> x
	highPlungeFrames = frames.InitAbilSlice(76)
	highPlungeFrames[action.ActionAttack] = 59
	highPlungeFrames[action.ActionSkill] = 60
	highPlungeFrames[action.ActionBurst] = 60
	highPlungeFrames[action.ActionSwap] = 64

	// low_plunge -> x
	lowPlungeFramesB = frames.InitAbilSlice(75)
	lowPlungeFramesB[action.ActionAttack] = 58
	lowPlungeFramesB[action.ActionSkill] = 58 // Assuming same as burst
	lowPlungeFramesB[action.ActionBurst] = 58
	lowPlungeFramesB[action.ActionWalk] = 74
	lowPlungeFramesB[action.ActionSwap] = 61

	// high_plunge -> x
	highPlungeFramesB = frames.InitAbilSlice(77)
	highPlungeFramesB[action.ActionAttack] = 60
	highPlungeFramesB[action.ActionSkill] = 59 // Assuming same as burst
	highPlungeFramesB[action.ActionBurst] = 59
	highPlungeFramesB[action.ActionWalk] = 76
	highPlungeFramesB[action.ActionSwap] = 63
}

// Low Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *char) LowPlungeAttack(p map[string]int) (action.Info, error) {
	defer c.Core.Player.SetAirborne(player.Grounded)
	switch c.Core.Player.Airborne() {
	case player.AirborneXianyun:
		if c.StatusIsActive(paramitaBuff) {
			return c.lowPlungeBXY(p), nil
		}
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
		CanQueueAfter:   lowPlungeFrames[action.ActionAttack],
		State:           action.PlungeAttackState,
	}
}

func (c *char) lowPlungeBXY(p map[string]int) action.Info {
	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not collision hit
	}

	if collision > 0 {
		c.plungeCollision(collisionHitmarkB)
	}

	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Low Plunge (E)",
		AttackTag:      attacks.AttackTagPlunge,
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeBlunt,
		PoiseDMG:       lowPlungePoiseDMG,
		Element:        attributes.Pyro,
		IgnoreInfusion: true,
		Durability:     25,
		Mult:           lowPlunge[c.TalentLvlAttack()],
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
	}
}

// High Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *char) HighPlungeAttack(p map[string]int) (action.Info, error) {
	defer c.Core.Player.SetAirborne(player.Grounded)
	switch c.Core.Player.Airborne() {
	case player.AirborneXianyun:
		if c.StatusIsActive(paramitaBuff) {
			return c.highPlungeBXY(p), nil
		}
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
	}
}

func (c *char) highPlungeBXY(p map[string]int) action.Info {
	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not collision hit
	}

	if collision > 0 {
		c.plungeCollision(collisionHitmarkB)
	}

	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "High Plunge (E)",
		AttackTag:      attacks.AttackTagPlunge,
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeBlunt,
		PoiseDMG:       highPlungePoiseDMG,
		Element:        attributes.Pyro,
		IgnoreInfusion: true,
		Durability:     25,
		Mult:           highPlunge[c.TalentLvlAttack()],
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
		CanQueueAfter:   highPlungeFramesB[action.ActionBurst],
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

	if c.StatusIsActive(paramitaBuff) {
		ai.Element = attributes.Pyro
		ai.IgnoreInfusion = true
	}

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, 1), delay, delay)
}
