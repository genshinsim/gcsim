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

// TODO: update low_plunge frames and hitboxes

const highPlungeHitmark = 37
const fieryHighPlungeHitmark = 41
const collisionHitmark = highPlungeHitmark - 6

const lowPlungeRadius = 3.0
const highPlungeRadius = 5.0

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
		FlatDmg:        c.c4FlatBonus(),
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
	ai.Mult += c.a1PlungeBuff()
	c.getApexDrive()

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, highPlungeRadius),
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

// Plunge normal falling attack damage queue generator
// Standard - Always part of high/low plunge attacks
func (c *char) plungeCollision(delay int) {
	ai := combat.AttackInfo{
		ActorIndex:   c.Index,
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
