package varesa

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

var (
	highPlungeFrames      []int
	fieryHighPlungeFrames []int
	lowPlungeFrames       []int
)

// based on gaming frames
// TODO: update frames

const lowPlungeHitmark = 39
const highPlungeHitmark = 39 + 3
const collisionHitmark = lowPlungeHitmark - 6

const lowPlungeRadius = 3.0
const highPlungeRadius = 5.0

const apexState = "apex-drive"

func init() {
	// low_plunge -> x
	lowPlungeFrames = frames.InitAbilSlice(84)
	lowPlungeFrames[action.ActionAttack] = 56
	lowPlungeFrames[action.ActionSkill] = 56
	lowPlungeFrames[action.ActionBurst] = 56
	lowPlungeFrames[action.ActionDash] = lowPlungeHitmark
	lowPlungeFrames[action.ActionSwap] = 67

	// high_plunge -> x
	highPlungeFrames = frames.InitAbilSlice(72) // Plunge -> Walk
	highPlungeFrames[action.ActionAttack] = 40
	highPlungeFrames[action.ActionSkill] = 40
	highPlungeFrames[action.ActionBurst] = 40
	highPlungeFrames[action.ActionDash] = 40
	highPlungeFrames[action.ActionJump] = 51
	highPlungeFrames[action.ActionSwap] = 37

	// fiery high_plunge -> x
	fieryHighPlungeFrames = frames.InitAbilSlice(90) // Plunge -> Walk
	fieryHighPlungeFrames[action.ActionAttack] = 47
	fieryHighPlungeFrames[action.ActionSkill] = 47
	fieryHighPlungeFrames[action.ActionBurst] = 47
	fieryHighPlungeFrames[action.ActionDash] = 40
	fieryHighPlungeFrames[action.ActionJump] = 79
	fieryHighPlungeFrames[action.ActionSwap] = 45
}

// Low Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *char) LowPlungeAttack(p map[string]int) (action.Info, error) {
	if c.Core.Player.CurrentState() == action.ChargeAttackState {
		return c.lowPlungeXY(p), nil
	}

	defer c.Core.Player.SetAirborne(player.Grounded)
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
		ActorIndex:     c.Index,
		Abil:           "Low Plunge",
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		AttackTag:      attacks.AttackTagPlunge,
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Electro,
		Durability:     25,
		Mult:           lowPlunge[c.TalentLvlAttack()],
	}

	cb := c.generatePlungeNightsoul
	c.exitNS = false
	if c.nightsoulState.HasBlessing() {
		ai.Abil = "Fiery Passion Low Plunge"
		ai.Mult = fieryLowPlunge[c.TalentLvlAttack()]
		cb = c.nightsoulState.ClearPoints
		c.exitNS = true
	}
	ai.Mult += c.a1PlungeBuff()
	c.getApexDrive()

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, lowPlungeRadius),
		lowPlungeHitmark,
		lowPlungeHitmark,
		c.a1Cancel,
		c.c2CB(),
	)
	c.Core.Tasks.Add(cb, lowPlungeHitmark)

	return action.Info{
		Frames:          frames.NewAbilFunc(lowPlungeFrames),
		AnimationLength: lowPlungeFrames[action.InvalidAction],
		CanQueueAfter:   lowPlungeFrames[action.ActionDash],
		State:           action.PlungeAttackState,
		OnRemoved:       c.clearNightsoul,
	}
}

// High Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *char) HighPlungeAttack(p map[string]int) (action.Info, error) {
	if c.Core.Player.CurrentState() == action.ChargeAttackState {
		return c.highPlungeXY(p), nil
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

	if collision > 0 {
		c.plungeCollision(collisionHitmark)
	}

	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "High Plunge",
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		AttackTag:      attacks.AttackTagPlunge,
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Electro,
		Durability:     25,
		Mult:           highPlunge[c.TalentLvlAttack()],
	}

	cb := c.generatePlungeNightsoul
	c.exitNS = false
	plungeFrames := highPlungeFrames
	if c.nightsoulState.HasBlessing() {
		ai.Abil = "Fiery Passion High Plunge"
		ai.Mult = fieryHighPlunge[c.TalentLvlAttack()]
		cb = c.nightsoulState.ClearPoints
		c.exitNS = true
		plungeFrames = fieryHighPlungeFrames
	}
	ai.Mult += c.a1PlungeBuff()
	c.getApexDrive()

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, highPlungeRadius),
		highPlungeHitmark,
		highPlungeHitmark,
		c.a1Cancel,
		c.c2CB(),
	)
	c.Core.Tasks.Add(cb, highPlungeHitmark)

	return action.Info{
		Frames:          frames.NewAbilFunc(plungeFrames),
		AnimationLength: plungeFrames[action.InvalidAction],
		CanQueueAfter:   plungeFrames[action.ActionDash],
		State:           action.PlungeAttackState,
		OnRemoved:       c.clearNightsoul,
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
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 0,
		Mult:       collision[c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 1), delay, delay)
}

func (c *char) getApexDrive() {
	if c.Base.Cons < 2 && !c.nightsoulState.HasBlessing() {
		return
	}
	c.AddStatus(apexState, 2*60, true) // TODO: duration?
	if c.Base.Cons >= 6 {
		c.AddEnergy("varesa-c6", 30)
	}
}
