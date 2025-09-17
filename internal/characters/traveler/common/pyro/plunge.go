package pyro

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var (
	highPlungeFrames [][]int
	lowPlungeFrames  [][]int
)

const (
	lowPlungeHitmark  = 47
	highPlungeHitmark = 48
	collisionHitmark  = lowPlungeHitmark - 6
)

const (
	lowPlungePoiseDMG = 100.0
	lowPlungeRadius   = 3.0
)

const (
	highPlungePoiseDMG = 150.0
	highPlungeRadius   = 5.0
)

func init() {
	// low_plunge -> x
	lowPlungeFrames = make([][]int, 2)

	lowPlungeFrames[0] = frames.InitAbilSlice(5000)

	lowPlungeFrames[1] = frames.InitAbilSlice(74)
	lowPlungeFrames[1][action.ActionAttack] = 58
	lowPlungeFrames[1][action.ActionSkill] = 58
	lowPlungeFrames[1][action.ActionBurst] = 58
	lowPlungeFrames[1][action.ActionDash] = lowPlungeHitmark
	lowPlungeFrames[1][action.ActionSwap] = 61

	// high_plunge -> x
	highPlungeFrames = make([][]int, 2)

	highPlungeFrames[0] = frames.InitAbilSlice(5000)

	highPlungeFrames[1] = frames.InitAbilSlice(75)
	highPlungeFrames[1][action.ActionAttack] = 58
	highPlungeFrames[1][action.ActionSkill] = 60
	highPlungeFrames[1][action.ActionBurst] = 59
	highPlungeFrames[1][action.ActionDash] = highPlungeHitmark
	highPlungeFrames[1][action.ActionSwap] = 62
}

// Low Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *Traveler) LowPlungeAttack(p map[string]int) (action.Info, error) {
	if c.gender == 0 {
		// aether not implemented
		return action.Info{}, fmt.Errorf("%v: action low_plunge not implemented", c.Base.Key)
	}
	defer c.Core.Player.SetAirborne(player.Grounded)
	switch c.Core.Player.Airborne() {
	case player.AirborneXianyun:
		return c.lowPlungeXY(p), nil
	default:
		return action.Info{}, errors.New("low_plunge can only be used while airborne")
	}
}

func (c *Traveler) lowPlungeXY(p map[string]int) action.Info {
	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not collision hit
	}

	if collision > 0 {
		c.plungeCollision(collisionHitmark)
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
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
	if c.Base.Cons >= 6 && c.nightsoulState.HasBlessing() {
		ai.Element = attributes.Pyro
		ai.IgnoreInfusion = true
		ai.AdditionalTags = []attacks.AdditionalTag{attacks.AdditionalTagNightsoul}
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: 1}, lowPlungeRadius),
		lowPlungeHitmark,
		lowPlungeHitmark,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(lowPlungeFrames[c.gender]),
		AnimationLength: lowPlungeFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   lowPlungeFrames[c.gender][action.ActionDash],
		State:           action.PlungeAttackState,
	}
}

// High Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *Traveler) HighPlungeAttack(p map[string]int) (action.Info, error) {
	if c.gender == 0 {
		// aether not implemented
		return action.Info{}, fmt.Errorf("%v: action low_plunge not implemented", c.Base.Key)
	}

	defer c.Core.Player.SetAirborne(player.Grounded)
	switch c.Core.Player.Airborne() {
	case player.AirborneXianyun:
		return c.highPlungeXY(p), nil
	default:
		return action.Info{}, errors.New("high_plunge can only be used while airborne")
	}
}

func (c *Traveler) highPlungeXY(p map[string]int) action.Info {
	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not collision hit
	}

	if collision > 0 {
		c.plungeCollision(collisionHitmark)
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
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
	if c.Base.Cons >= 6 && c.nightsoulState.HasBlessing() {
		ai.Element = attributes.Pyro
		ai.IgnoreInfusion = true
		ai.AdditionalTags = []attacks.AdditionalTag{attacks.AdditionalTagNightsoul}
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: 1}, highPlungeRadius),
		highPlungeHitmark,
		highPlungeHitmark,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(highPlungeFrames[c.gender]),
		AnimationLength: highPlungeFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   highPlungeFrames[c.gender][action.ActionDash],
		State:           action.PlungeAttackState,
	}
}

// Plunge normal falling attack damage queue generator
// Standard - Always part of high/low plunge attacks
func (c *Traveler) plungeCollision(delay int) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Plunge Collision",
		AttackTag:  attacks.AttackTagPlunge,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Physical,
		Durability: 0,
		Mult:       collision[c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: 1}, 1), delay, delay)
}
