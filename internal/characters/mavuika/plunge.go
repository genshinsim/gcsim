package mavuika

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
var lowBikePlungeFrames []int
var highBikePlungeFrames []int

const lowPlungeHitmark = 37
const highPlungeHitmark = 41
const lowBikePungeHitmark = 41
const highBikePungeHitmark = 45
const collisionHitmark = lowPlungeHitmark - 6

const lowPlungePoiseDMG = 150.0
const lowPlungeRadius = 3.0

const highPlungePoiseDMG = 200.0
const highPlungeRadius = 5.0

const bikePlungePoiseDMG = 150.0
const bikePlungeRadius = 5.0

func init() {
	// low_plunge -> x
	lowPlungeFrames = frames.InitAbilSlice(80) // low_plunge -> Jump
	lowPlungeFrames[action.ActionAttack] = 51
	lowPlungeFrames[action.ActionCharge] = 52
	lowPlungeFrames[action.ActionSkill] = 37 // low_plunge -> skill[recast=1]
	lowPlungeFrames[action.ActionBurst] = 51
	lowPlungeFrames[action.ActionDash] = lowPlungeHitmark
	lowPlungeFrames[action.ActionWalk] = 79
	lowPlungeFrames[action.ActionSwap] = 62

	// high_plunge -> x
	highPlungeFrames = frames.InitAbilSlice(83)
	highPlungeFrames[action.ActionAttack] = 54
	highPlungeFrames[action.ActionCharge] = 55
	highPlungeFrames[action.ActionSkill] = 40 // low_plunge -> skill[recast=1]
	highPlungeFrames[action.ActionBurst] = 53
	highPlungeFrames[action.ActionDash] = highPlungeHitmark
	highPlungeFrames[action.ActionWalk] = 82
	highPlungeFrames[action.ActionSwap] = 65

	// Flamestrider low_plunge -> X
	lowBikePlungeFrames = frames.InitAbilSlice(77) // low_plunge -> Walk
	lowBikePlungeFrames[action.ActionAttack] = 60
	lowBikePlungeFrames[action.ActionCharge] = 60
	lowBikePlungeFrames[action.ActionSkill] = 41 // low_plunge -> skill[recast=1]
	lowBikePlungeFrames[action.ActionBurst] = 61
	lowBikePlungeFrames[action.ActionDash] = lowBikePungeHitmark
	lowBikePlungeFrames[action.ActionJump] = 76
	lowBikePlungeFrames[action.ActionSwap] = 75

	// Flamestrider high_plunge -> X
	highBikePlungeFrames = frames.InitAbilSlice(80) // low_plunge -> Walk
	highBikePlungeFrames[action.ActionAttack] = 63
	highBikePlungeFrames[action.ActionCharge] = 63
	highBikePlungeFrames[action.ActionSkill] = 44 // low_plunge -> skill[recast=1]
	highBikePlungeFrames[action.ActionBurst] = 64
	highBikePlungeFrames[action.ActionDash] = highBikePungeHitmark
	highBikePlungeFrames[action.ActionJump] = 79
	highBikePlungeFrames[action.ActionSwap] = 78
}

// Low Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *char) LowPlungeAttack(p map[string]int) (action.Info, error) {
	defer c.Core.Player.SetAirborne(player.Grounded)
	if c.Core.Player.Airborne() == player.AirborneXianyun || c.canBikePlunge {
		if c.nightsoulState.HasBlessing() {
			return c.bikePlungeAttack(lowBikePlungeFrames, lowPlungeHitmark), nil
		}
		return c.lowPlungeXY(p), nil
	}
	return action.Info{}, errors.New("low_plunge can only be used while airborne")
}

// Also used for low plunge from walked bike jump -> NS expired -> plunge
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

// High Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *char) HighPlungeAttack(p map[string]int) (action.Info, error) {
	defer c.Core.Player.SetAirborne(player.Grounded)
	switch c.Core.Player.Airborne() {
	case player.AirborneXianyun:
		if c.nightsoulState.HasBlessing() && c.armamentState == bike {
			c.bikePlungeAttack(highBikePlungeFrames, highPlungeHitmark)
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
		CanQueueAfter:   highPlungeFrames[action.ActionDash],
		State:           action.PlungeAttackState,
	}
}

// Flamestrider plunge attack damage queue generator
func (c *char) bikePlungeAttack(bikePlungeFrames []int, delay int) action.Info {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Flamestrider Plunge",
		AttackTag:      attacks.AttackTagPlunge,
		ICDTag:         attacks.ICDTagMavuikaFlamestrider,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeBlunt,
		PoiseDMG:       bikePlungePoiseDMG,
		Element:        attributes.Pyro,
		Durability:     25,
		Mult:           skillPlunge[c.TalentLvlSkill()],
		HitlagFactor:   0.1,
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, bikePlungeRadius),
		delay,
		delay,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(bikePlungeFrames),
		AnimationLength: bikePlungeFrames[action.InvalidAction],
		CanQueueAfter:   bikePlungeFrames[action.ActionDash],
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
