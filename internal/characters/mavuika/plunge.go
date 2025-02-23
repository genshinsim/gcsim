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
var bikePlungeFrames []int

const lowPlungeHitmark = 37
const highPlungeHitmark = 41
const bikePlungeHitmark = 41
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
	bikePlungeFrames = frames.InitAbilSlice(77) // low_plunge -> Walk
	bikePlungeFrames[action.ActionAttack] = 60
	bikePlungeFrames[action.ActionCharge] = 60
	bikePlungeFrames[action.ActionSkill] = 41 // low_plunge -> skill[recast=1]
	bikePlungeFrames[action.ActionBurst] = 61
	bikePlungeFrames[action.ActionDash] = bikePlungeHitmark
	bikePlungeFrames[action.ActionJump] = 76
	bikePlungeFrames[action.ActionSwap] = 75
}

// Low Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *char) LowPlungeAttack(p map[string]int) (action.Info, error) {
	defer c.Core.Player.SetAirborne(player.Grounded)
	switch c.Core.Player.Airborne() {
	case player.AirborneXianyun:
		return c.lowPlungeXY(p), nil
	default:
		if c.canBikePlunge {
			c.bikePlungeAttack()
		}
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

// Flamestrider falling attack damage queue generator
func (c *char) bikePlungeAttack() action.Info {
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
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, bikePlungeRadius),
		bikePlungeHitmark,
		bikePlungeHitmark,
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
