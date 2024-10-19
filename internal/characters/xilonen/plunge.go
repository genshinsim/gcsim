package xilonen

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

var lowPlungeFrames []int
var highPlungeFrames []int
var skillHighPlungeFrames []int

const lowPlungeHitmark = 47
const highPlungeHitmark = 50
const collisionHitmark = lowPlungeHitmark - 6

const skillHighPlungeHitmark = 39
const skillCollisionHitmark = skillHighPlungeHitmark - 6

const lowPlungePoiseDMG = 100.0
const lowPlungeRadius = 3.0

const highPlungePoiseDMG = 150.0
const highPlungeRadius = 5.0

func init() {
	// low_plunge -> x
	lowPlungeFrames = frames.InitAbilSlice(75)
	lowPlungeFrames[action.ActionAttack] = 59
	lowPlungeFrames[action.ActionSkill] = 59
	lowPlungeFrames[action.ActionBurst] = 58
	lowPlungeFrames[action.ActionDash] = lowPlungeHitmark
	lowPlungeFrames[action.ActionSwap] = 60
	lowPlungeFrames[action.ActionWalk] = 74

	// high_plunge -> x
	highPlungeFrames = frames.InitAbilSlice(75)
	highPlungeFrames[action.ActionAttack] = 62
	highPlungeFrames[action.ActionSkill] = 60
	highPlungeFrames[action.ActionBurst] = 61
	highPlungeFrames[action.ActionDash] = highPlungeHitmark
	highPlungeFrames[action.ActionJump] = 74
	highPlungeFrames[action.ActionSwap] = 61

	skillHighPlungeFrames = frames.InitAbilSlice(77)
	skillHighPlungeFrames[action.ActionAttack] = 55
	skillHighPlungeFrames[action.ActionSkill] = 56
	skillHighPlungeFrames[action.ActionBurst] = 55
	skillHighPlungeFrames[action.ActionDash] = skillHighPlungeHitmark
	skillHighPlungeFrames[action.ActionSwap] = 57
	skillHighPlungeFrames[action.ActionWalk] = 66
}

// Low Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *char) LowPlungeAttack(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		c.c6()
	}
	defer c.Core.Player.SetAirborne(player.Grounded)
	switch c.Core.Player.Airborne() {
	case player.AirborneXianyun:
		if c.nightsoulState.HasBlessing() {
			return action.Info{}, errors.New("xilonen cannot low_plunge while in nightsoul blessing")
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
		UseDef:     true,
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, lowPlungeRadius),
		lowPlungeHitmark,
		lowPlungeHitmark,
		c.a1cb,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(lowPlungeFrames),
		AnimationLength: lowPlungeFrames[action.InvalidAction],
		CanQueueAfter:   lowPlungeFrames[action.ActionDash],
		State:           action.PlungeAttackState,
	}
}

// High Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *char) HighPlungeAttack(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		c.c6()
	}
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
		UseDef:     true,
	}
	highPlungeFrames := highPlungeFrames
	collisionHitmark := collisionHitmark
	if c.nightsoulState.HasBlessing() {
		ai.Element = attributes.Geo
		ai.IgnoreInfusion = true
		ai.AdditionalTags = []attacks.AdditionalTag{attacks.AdditionalTagNightsoul}
		ai.Mult += c.c6DmgMult()

		highPlungeFrames = skillHighPlungeFrames
		collisionHitmark = skillCollisionHitmark
	}

	if collision > 0 {
		c.plungeCollision(collisionHitmark)
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, highPlungeRadius),
		highPlungeHitmark,
		highPlungeHitmark,
		c.a1cb,
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
		UseDef:     true,
	}
	if c.nightsoulState.HasBlessing() {
		ai.Element = attributes.Geo
		ai.IgnoreInfusion = true
		ai.AdditionalTags = []attacks.AdditionalTag{attacks.AdditionalTagNightsoul}
		ai.Mult += c.c6DmgMult()
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, 1), delay, delay)
}
