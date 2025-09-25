package varesa

import (
	"errors"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var (
	highPlungeFrames      []int
	fieryHighPlungeFrames []int

	xianyunLowPlungeFrames  []int
	xianyunHighPlungeFrames []int

	xianyunFieryLowPlungeFrames  []int
	xianyunFieryHighPlungeFrames []int
)

// TODO: low_plunge

const (
	highPlungeHitmark      = 37
	fieryHighPlungeHitmark = 41
	collisionHitmark       = highPlungeHitmark - 6
)

const (
	xianyunLowPlungeHitmark  = 37
	xianyunHighPlungeHitmark = 38
)

const (
	xianyunLowPlungeNonNSWalk  = 72
	xianyunHighPlungeNonNSWalk = 72
)

const (
	xianyunFieryLowPlungeHitmark  = 40
	xianyunFieryHighPlungeHitmark = 40
)

const (
	lowPlungeRadius  = 3.0
	highPlungeRadius = 5.0
)

const (
	apexState    = "apex-drive"
	apexDuration = 140
)

func init() {
	// high_plunge -> x
	highPlungeFrames = frames.InitAbilSlice(72) // Plunge -> Walk
	highPlungeFrames[action.ActionAttack] = 40
	highPlungeFrames[action.ActionCharge] = 47
	highPlungeFrames[action.ActionSkill] = 40
	highPlungeFrames[action.ActionBurst] = 40
	highPlungeFrames[action.ActionDash] = 40
	highPlungeFrames[action.ActionJump] = 51
	highPlungeFrames[action.ActionSwap] = 37

	// fiery high_plunge -> x
	fieryHighPlungeFrames = frames.InitAbilSlice(90) // Plunge -> Walk
	fieryHighPlungeFrames[action.ActionAttack] = 47
	fieryHighPlungeFrames[action.ActionCharge] = 46
	fieryHighPlungeFrames[action.ActionSkill] = 47
	fieryHighPlungeFrames[action.ActionBurst] = 47
	fieryHighPlungeFrames[action.ActionDash] = 40
	fieryHighPlungeFrames[action.ActionJump] = 79
	fieryHighPlungeFrames[action.ActionSwap] = 45

	// xianyun low_plunge -> x
	xianyunLowPlungeFrames = frames.InitAbilSlice(75) // Plunge -> Walk
	xianyunLowPlungeFrames[action.ActionAttack] = 40
	xianyunLowPlungeFrames[action.ActionCharge] = 39
	xianyunLowPlungeFrames[action.ActionSkill] = 39
	xianyunLowPlungeFrames[action.ActionBurst] = 39
	xianyunLowPlungeFrames[action.ActionDash] = xianyunLowPlungeHitmark
	xianyunLowPlungeFrames[action.ActionJump] = 50
	xianyunLowPlungeFrames[action.ActionSwap] = 38

	// xianyun high_plunge -> x
	xianyunHighPlungeFrames = frames.InitAbilSlice(77) // Plunge -> Walk
	xianyunHighPlungeFrames[action.ActionAttack] = 39
	xianyunHighPlungeFrames[action.ActionCharge] = 40
	xianyunHighPlungeFrames[action.ActionSkill] = 39
	xianyunHighPlungeFrames[action.ActionBurst] = 39
	xianyunHighPlungeFrames[action.ActionDash] = xianyunHighPlungeHitmark
	xianyunHighPlungeFrames[action.ActionJump] = 49
	xianyunHighPlungeFrames[action.ActionSwap] = 39

	// xianyun nightsoul low_plunge -> x
	xianyunFieryLowPlungeFrames = frames.InitAbilSlice(91) // Plunge -> Walk
	xianyunFieryLowPlungeFrames[action.ActionAttack] = 46
	xianyunFieryLowPlungeFrames[action.ActionCharge] = 45
	xianyunFieryLowPlungeFrames[action.ActionSkill] = 45
	xianyunFieryLowPlungeFrames[action.ActionBurst] = 46
	xianyunFieryLowPlungeFrames[action.ActionDash] = xianyunFieryLowPlungeHitmark
	xianyunFieryLowPlungeFrames[action.ActionJump] = 76
	xianyunFieryLowPlungeFrames[action.ActionSwap] = 46

	// xianyun nightsoul high_plunge -> x
	xianyunFieryHighPlungeFrames = frames.InitAbilSlice(91) // Plunge -> Walk
	xianyunFieryHighPlungeFrames[action.ActionAttack] = 47
	xianyunFieryHighPlungeFrames[action.ActionCharge] = 47
	xianyunFieryHighPlungeFrames[action.ActionSkill] = 47
	xianyunFieryHighPlungeFrames[action.ActionBurst] = 47
	xianyunFieryHighPlungeFrames[action.ActionDash] = xianyunFieryHighPlungeHitmark
	xianyunFieryHighPlungeFrames[action.ActionJump] = 77
	xianyunFieryHighPlungeFrames[action.ActionSwap] = 46
}

// High Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *char) HighPlungeAttack(p map[string]int) (action.Info, error) {
	defer c.Core.Player.SetAirborne(player.Grounded)
	if c.Core.Player.CurrentState() == action.ChargeAttackState {
		return c.highPlungeCA(p), nil
	}
	switch c.Core.Player.Airborne() {
	case player.AirborneXianyun:
		return c.highPlungeXY(p), nil
	default:
		return action.Info{}, errors.New("high_plunge can only be used after charge or while airborne")
	}
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
		return action.Info{}, errors.New("low_plunge can only be used while airborne")
	}
}

func (c *char) highPlungeCA(p map[string]int) action.Info {
	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not collision hit
	}

	if collision > 0 {
		c.plungeCollision(collisionHitmark)
	}

	ai := info.AttackInfo{
		ActorIndex:     c.Index(),
		Abil:           "High Plunge",
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		AttackTag:      attacks.AttackTagPlunge,
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Electro,
		Durability:     25,
		Mult:           highPlunge[c.TalentLvlAttack()],
		FlatDmg:        c.a1PlungeBonus() + c.c4FlatBonus(),
		HitlagFactor:   0.1,
	}

	hitmark := highPlungeHitmark
	c.exitNS = false
	plungeFrames := highPlungeFrames
	if c.nightsoulState.HasBlessing() {
		ai.Abil = "Fiery Passion High Plunge"
		ai.Mult = fieryHighPlunge[c.TalentLvlAttack()]
		hitmark = fieryHighPlungeHitmark
		c.exitNS = true
		plungeFrames = fieryHighPlungeFrames

		c.QueueCharTask(c.nightsoulState.ClearPoints, 2)
	} else {
		c.QueueCharTask(c.generatePlungeNightsoul, hitmark)
	}
	c.getApexDrive()

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: 1}, highPlungeRadius),
		hitmark,
		hitmark,
		c.a1Cancel,
		c.c2CB(),
		c.c4CB,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(plungeFrames),
		AnimationLength: plungeFrames[action.InvalidAction],
		CanQueueAfter:   plungeFrames[action.ActionDash],
		State:           action.PlungeAttackState,
		OnRemoved:       c.clearNightsoulCB,
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

	ai := info.AttackInfo{
		ActorIndex:     c.Index(),
		Abil:           "High Plunge",
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		AttackTag:      attacks.AttackTagPlunge,
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Electro,
		Durability:     25,
		Mult:           highPlunge[c.TalentLvlAttack()],
		FlatDmg:        c.a1PlungeBonus() + c.c4FlatBonus(),
		HitlagFactor:   0.1,
	}

	hitmark := xianyunHighPlungeHitmark
	c.exitNS = false
	plungeFrames := c.newAbilFuncXYPlunge(xianyunHighPlungeFrames, xianyunHighPlungeNonNSWalk)
	if c.nightsoulState.HasBlessing() {
		ai.Abil = "Fiery Passion High Plunge"
		ai.Mult = fieryHighPlunge[c.TalentLvlAttack()]
		hitmark = xianyunFieryHighPlungeHitmark
		c.exitNS = true
		plungeFrames = frames.NewAbilFunc(xianyunFieryHighPlungeFrames)

		c.QueueCharTask(c.nightsoulState.ClearPoints, 2)
	} else {
		c.QueueCharTask(c.generatePlungeNightsoul, hitmark)
	}
	c.getApexDrive()

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: 1}, highPlungeRadius),
		hitmark,
		hitmark,
		c.a1Cancel,
		c.c2CB(),
		c.c4CB,
	)

	return action.Info{
		Frames:          plungeFrames,
		AnimationLength: plungeFrames(action.InvalidAction),
		CanQueueAfter:   plungeFrames(action.ActionDash),
		State:           action.PlungeAttackState,
		OnRemoved:       c.clearNightsoulCB,
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

	ai := info.AttackInfo{
		ActorIndex:     c.Index(),
		Abil:           "Low Plunge",
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		AttackTag:      attacks.AttackTagPlunge,
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Electro,
		Durability:     25,
		Mult:           lowPlunge[c.TalentLvlAttack()],
		FlatDmg:        c.a1PlungeBonus() + c.c4FlatBonus(),
		HitlagFactor:   0.1,
	}

	hitmark := xianyunLowPlungeHitmark
	c.exitNS = false
	plungeFrames := c.newAbilFuncXYPlunge(xianyunLowPlungeFrames, xianyunLowPlungeNonNSWalk)
	if c.nightsoulState.HasBlessing() {
		ai.Abil = "Fiery Passion Low Plunge"
		ai.Mult = fieryLowPlunge[c.TalentLvlAttack()]
		hitmark = xianyunFieryLowPlungeHitmark
		c.exitNS = true
		plungeFrames = frames.NewAbilFunc(xianyunFieryLowPlungeFrames)

		c.QueueCharTask(c.nightsoulState.ClearPoints, 2)
	} else {
		c.QueueCharTask(c.generatePlungeNightsoul, hitmark)
	}
	c.getApexDrive()

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: 1}, lowPlungeRadius),
		hitmark,
		hitmark,
		c.a1Cancel,
		c.c2CB(),
		c.c4CB,
	)

	return action.Info{
		Frames:          plungeFrames,
		AnimationLength: plungeFrames(action.InvalidAction),
		CanQueueAfter:   plungeFrames(action.ActionDash),
		State:           action.PlungeAttackState,
		OnRemoved:       c.clearNightsoulCB,
	}
}

func (c *char) newAbilFuncXYPlunge(slice []int, nonNSWalk int) func(action.Action) int {
	// This is different because the walk frames after plunge change based on if in nightsoul or not
	return func(next action.Action) int {
		if !c.nightsoulState.HasBlessing() && next == action.ActionWalk {
			return nonNSWalk
		}
		return slice[next]
	}
}

// Plunge normal falling attack damage queue generator
// Standard - Always part of high/low plunge attacks
func (c *char) plungeCollision(delay int) {
	ai := info.AttackInfo{
		ActorIndex:   c.Index(),
		Abil:         "Plunge Collision",
		AttackTag:    attacks.AttackTagPlunge,
		ICDTag:       attacks.ICDTagNone,
		ICDGroup:     attacks.ICDGroupDefault,
		StrikeType:   attacks.StrikeTypeDefault,
		Element:      attributes.Electro,
		Durability:   0,
		Mult:         collision[c.TalentLvlAttack()],
		HitlagFactor: 0.02,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 1), delay, delay)
}

func (c *char) getApexDrive() {
	if c.Base.Cons < 2 && !c.nightsoulState.HasBlessing() {
		return
	}
	c.AddStatus(apexState, apexDuration, true)
	if c.Base.Cons >= 6 {
		c.AddEnergy("varesa-c6", 30)
	}
}
